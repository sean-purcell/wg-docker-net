package wg

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"

	"github.com/docker/go-plugins-helpers/network"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"gopkg.in/go-ini/ini.v1"
)

type WgConfig struct {
	Path       string
	ListenPort uint
	Net        *net.IPNet
	PeerNets   []*net.IPNet
}

func ParseWgConfig(path string) (*WgConfig, error) {
	ini, err := ini.LoadSources(ini.LoadOptions{AllowNonUniqueSections: true}, path)
	if err != nil {
		return nil, err
	}

	intf := ini.Section("Interface")
	if intf == nil {
		return nil, fmt.Errorf("[Interface] section missing in config: %s\n", path)
	}

	key, err := intf.GetKey("ListenPort")
	if err != nil {
		return nil, err
	}
	ListenPort, err := key.Uint()
	if err != nil {
		return nil, err
	}

	key, err = intf.GetKey("Address")
	if err != nil {
		return nil, err
	}

	ip, Net, err := net.ParseCIDR(key.String())
	Net.IP = ip
	if err != nil {
		return nil, err
	}

	sections, err := ini.SectionsByName("Peer")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Num sections: %v\n", len(sections))
	PeerNets := make([]*net.IPNet, 0)
	for _, section := range sections {
		key, err := section.GetKey("AllowedIPs")
		if err != nil {
			return nil, err
		}
		fmt.Printf("AllowedIps: %s\n", key.Value())
		for _, addr := range key.Strings(",") {
			_, peerNet, err := net.ParseCIDR(strings.TrimSpace(addr))
			if err != nil {
				return nil, err
			}
			PeerNets = append(PeerNets, peerNet)
		}
	}

	Path := path
	return &WgConfig{Path, ListenPort, Net, PeerNets}, nil
}

func (t *WgConfig) StartInterface(nl *netlink.Handle) (netlink.Link, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	currentNs, err := netns.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		netns.Set(currentNs)
	}()

	log.Printf("Bringing up wireguard interface at %s\n", t.Path)
	cmd := exec.Command("wg-quick", "up", t.Path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command failed: %v\n", err)
		return nil, err
	}
	log.Printf("Output: %s\n", string(output))

	links, err := nl.LinkList()
	if err != nil {
		return nil, fmt.Errorf("Failed to enumerate links %v", err)
	}

	for _, link := range links {
		if link.Type() == "wireguard" {
			return link, nil
		}
	}

	return nil, fmt.Errorf("Wireguard interface not found")
}

func (t *WgConfig) GetRoutes(gateway net.IP) []*network.StaticRoute {
	routes := make([]*network.StaticRoute, len(t.PeerNets))
	for i, peer := range t.PeerNets {
		routes[i] = &network.StaticRoute{
			Destination: peer.String(),
			RouteType:   0,
			NextHop:     gateway.String(),
		}
	}
	return routes
}
