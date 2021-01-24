package wg

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/vishvananda/netns"
	"gopkg.in/go-ini/ini.v1"
)

type WgConfig struct {
	path          string
	listenPort    uint
	address       string
	peerAddresses []string
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
	listenPort, err := key.Uint()
	if err != nil {
		return nil, err
	}

	key, err = intf.GetKey("Address")
	if err != nil {
		return nil, err
	}
	address := key.String()

	sections, err := ini.SectionsByName("Peer")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Num sections: %v\n", len(sections))
	peerAddresses := make([]string, 0)
	for _, section := range sections {
		key, err := section.GetKey("AllowedIPs")
		if err != nil {
			return nil, err
		}
		fmt.Printf("AllowedIps: %s\n", key.Value())
		for _, addr := range key.Strings(",") {
			peerAddresses = append(peerAddresses, strings.TrimSpace(addr))
		}
	}

	return &WgConfig{path, listenPort, address, peerAddresses}, nil
}

func (t *WgConfig) StartInterface() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	currentNs, err := netns.Get()
	if err != nil {
		return err
	}
	defer func() {
		netns.Set(currentNs)
	}()

	log.Printf("Bringing up wireguard interface at %s\n", t.path)
	cmd := exec.Command("wg-quick", "up", t.path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command failed: %v\n", err)
		return err
	}
	log.Printf("Output: %s\n", string(output))

	return nil
}
