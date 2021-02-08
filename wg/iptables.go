package wg

import (
    "github.com/coreos/go-iptables/iptables"
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
)

type Iptables struct {
    i *iptables.IPTables
}

func jumpRule(target string) []string {

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
    if err = i.i.AppendUnique(table, sourceChain, "-p", "udp", "-j", chain); err != nil {
        return err
    }
    return nil
}

func (i *Iptables) deleteChain(table, sourceChain, chain string) error {
    if err := i.i.Delete(table, sourceChain, "-p", "udp", "-j", chain); err != nil {
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

func (i *Iptables) Delete() error {
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

func (i *IpTables) SetupForwarding(source, endpoint net.IP, port int) error {
    if err := i.i.Insert(nat, source_pre, 0, "-j", "DNAT", "-p", "udp",
}

func (i *IpTables) RemoveForwarding(source, endpoint net.IP, port int) error {
}
