package wg

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/go-plugins-helpers/network"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type Network struct {
	ns   netns.NsHandle
	nl   *netlink.Handle
	name *string
	conf *WgConfig
}

func getOpt(options map[string]interface{}, name string) *string {
	val, ok := options[name]
	if ok {
		str := val.(string)
		return &str
	} else {
		return nil
	}
}

func CreateNetwork(data *network.IPAMData, options map[string]interface{}, rootNs netns.NsHandle) (*Network, error) {
	var ns netns.NsHandle
	var err error

	confPath := getOpt(options, "wg.wgconf")

	if confPath == nil {
		return nil, fmt.Errorf("Wireguard config file not present")
	}

	conf, err := ParseWgConfig(*confPath)
	if err != nil {
		return nil, err
	}
	str := spew.Sdump(*conf)
	log.Printf("Loaded wireguard config: %s\n", str)

	name := getOpt(options, "wg.namespace")
	if name != nil {
		log.Printf("Creating namespace: %s\n", *name)
		ns, err = netns.NewNamed(*name)
		if err != nil {
			return nil, err
		}
	} else {
		log.Printf("Creating anonymous namespace\n")
		ns, err = netns.New()
		if err != nil {
			return nil, err
		}
	}

	log.Printf("Created namespace at fd %d\n", ns)

	nl, err := netlink.NewHandleAt(ns)
	if err != nil {
		return nil, err
	}

	return &Network{
		ns,
		nl,
		name,
		conf,
	}, nil
}

func (t *Network) Delete() error {
	if t.name != nil {
		err := netns.DeleteNamed(*t.name)
		if err != nil {
			return err
		}
	}

	return nil
}
