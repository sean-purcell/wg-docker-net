package wg

import (
	"log"

	"github.com/docker/go-plugins-helpers/network"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type Network struct {
	ns netns.NsHandle
	nl *netlink.Handle
}

func nsName(options map[string]interface{}) *string {
	name, ok := options["wg.namespace"]
	if ok {
		nameStr := name.(string)
		return &nameStr
	} else {
		return nil
	}
}

func CreateNetwork(data *network.IPAMData, options map[string]interface{}) (*Network, error) {
	var ns netns.NsHandle
	var err error
	if name := nsName(options); name != nil {
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

	nl, err := netlink.NewHandleAt(ns)
	if err != nil {
		return nil, err
	}

	return &Network{
		ns,
		nl,
	}, nil
}
