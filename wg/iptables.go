package wg

import (
	"net"
	"runtime"
	"strconv"

	"github.com/coreos/go-iptables/iptables"
	"github.com/vishvananda/netns"
)

const (
	nat    = "nat"
	filter = "filter"

	chain_prefix = "WG-DOCKER-"

	source_pre     = "PREROUTING"
	source_post    = "POSTROUTING"
	source_forward = "FORWARD"

	pre     = chain_prefix + source_pre
	post    = chain_prefix + source_post
	forward = chain_prefix + source_forward

	jump  = "-j"
	proto = "-p"
	udp   = "udp"
)

type Iptables struct {
	i *iptables.IPTables
}

func jumpRule(target string) []string {
	return []string{jump, target, proto, udp}
}

func (i *Iptables) createOrClearChain(table, sourceChain, chain string) error {
	exists, err := i.i.ChainExists(table, chain)
	if err != nil {
		return err
	}

	if exists {
		err = i.i.ClearChain(table, chain)
	} else {
		err = i.i.NewChain(table, chain)
	}

	if err != nil {
		return err
	}

	if err = i.i.Append(table, chain, "-j", "RETURN"); err != nil {
		return err
	}
	if err = i.i.AppendUnique(table, sourceChain, jumpRule(chain)...); err != nil {
		return err
	}
	return nil
}

func (i *Iptables) deleteChain(table, sourceChain, chain string) error {
	if err := i.i.Delete(table, sourceChain, jumpRule(chain)...); err != nil {
		return err
	}
	if err := i.i.ClearAndDeleteChain(table, chain); err != nil {
		return err
	}
	return nil
}

func CreateIptables() (*Iptables, error) {
	iptables, err := iptables.New()
	if err != nil {
		return nil, err
	}
	i := &Iptables{iptables}

	if err = i.createOrClearChain(nat, source_pre, pre); err != nil {
		return nil, err
	}
	if err = i.createOrClearChain(nat, source_post, post); err != nil {
		return nil, err
	}
	if err = i.createOrClearChain(filter, source_forward, forward); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Iptables) Delete(ns netns.NsHandle) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	currentNs, err := netns.Get()
	if err != nil {
		return err
	}
	defer func() {
		netns.Set(currentNs)
		_ = currentNs.Close()
	}()

	err = netns.Set(ns)
	if err != nil {
		return err
	}

	if err := i.deleteChain(nat, source_pre, pre); err != nil {
		return err
	}
	if err := i.deleteChain(nat, source_post, post); err != nil {
		return err
	}
	if err := i.deleteChain(filter, source_forward, forward); err != nil {
		return err
	}
	return nil
}

func snatRule(source, endpoint net.IP, port uint) []string {
	return []string{jump, "SNAT", proto, udp, "--source", source.String(), "--source-port", strconv.Itoa(int(port)), "--to-source", endpoint.String()}
}

func dnatRule(source, endpoint net.IP, port uint) []string {
	return []string{jump, "DNAT", proto, udp, "--destination", endpoint.String(), "--destination-port", strconv.Itoa(int(port)), "--to-destination", source.String()}
}

func forwardOutRule(source net.IP, port uint) []string {
	return []string{jump, "ACCEPT", proto, udp, "--destination", source.String(), "--destination-port", strconv.Itoa(int(port))}
}

func forwardInRule(source net.IP, port uint) []string {
	return []string{jump, "ACCEPT", proto, udp, "--source", source.String(), "--source-port", strconv.Itoa(int(port))}
}

func (i *Iptables) SetupForwarding(ns netns.NsHandle, source, endpoint net.IP, port uint) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	currentNs, err := netns.Get()
	if err != nil {
		return err
	}
	defer func() {
		netns.Set(currentNs)
		_ = currentNs.Close()
	}()

	err = netns.Set(ns)
	if err != nil {
		return err
	}

	if err := i.i.Insert(nat, pre, 1, dnatRule(source, endpoint, port)...); err != nil {
		return err
	}
	if err := i.i.Insert(nat, post, 1, snatRule(source, endpoint, port)...); err != nil {
		return err
	}
	if err := i.i.Insert(filter, forward, 1, forwardOutRule(source, port)...); err != nil {
		return err
	}
	if err := i.i.Insert(filter, forward, 1, forwardInRule(source, port)...); err != nil {
		return err
	}
	return nil
}

func (i *Iptables) RemoveForwarding(ns netns.NsHandle, source, endpoint net.IP, port uint) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	currentNs, err := netns.Get()
	if err != nil {
		return err
	}
	defer func() {
		netns.Set(currentNs)
		_ = currentNs.Close()
	}()

	err = netns.Set(ns)
	if err != nil {
		return err
	}

	if err := i.i.DeleteIfExists(nat, pre, dnatRule(source, endpoint, port)...); err != nil {
		return err
	}
	if err := i.i.DeleteIfExists(nat, post, snatRule(source, endpoint, port)...); err != nil {
		return err
	}
	if err := i.i.DeleteIfExists(filter, forward, forwardOutRule(source, port)...); err != nil {
		return err
	}
	if err := i.i.DeleteIfExists(filter, forward, forwardInRule(source, port)...); err != nil {
		return err
	}
	return nil
}
