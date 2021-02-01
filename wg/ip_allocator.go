package wg

import (
	"encoding/binary"
	"fmt"
	"net"
)

type IpAllocator struct {
	usedAddresses map[uint32]struct{}
	nextAddress   uint32
	mask          uint32
	upperBound    uint32
}

func bytesToUint(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func uintToBytes(val uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, val)
	return bytes
}

func CreateIpAllocator(subnet *net.IPNet) *IpAllocator {
	mask := bytesToUint(subnet.Mask)
	nextAddress := (bytesToUint(subnet.IP.To4()) & mask) + 1
	upperBound := (nextAddress | (^mask)) + 1
	usedAddresses := make(map[uint32]struct{}, 0)
	return &IpAllocator{
		usedAddresses, nextAddress, mask, upperBound,
	}
}

func (t *IpAllocator) IsUsed(ip net.IP) bool {
	_, ok := t.usedAddresses[bytesToUint(ip.To4())]
	return ok
}

func (t *IpAllocator) MarkUsed(ip net.IP) {
	t.usedAddresses[bytesToUint(ip.To4())] = struct{}{}
}

func (t *IpAllocator) FindAddress() (*net.IPNet, error) {
	// TODO: This doesn't work if the subnet borders the top of the ipv4 address space
	for ; t.nextAddress < t.upperBound; t.nextAddress++ {
		if _, ok := t.usedAddresses[t.nextAddress]; !ok {
			t.usedAddresses[t.nextAddress] = struct{}{}
			ip := uintToBytes(t.nextAddress)
			t.nextAddress++
			return &net.IPNet{ip, uintToBytes(t.mask)}, nil
		}
	}
	return nil, fmt.Errorf("No unused addresses remaining")
}
