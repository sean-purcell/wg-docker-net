package wg

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/go-plugins-helpers/network"
	"github.com/vishvananda/netns"
)

type Driver struct {
	networks map[string]*Network
	rootNs   netns.NsHandle
}

func notSupported(method string) error {
	return fmt.Errorf("[%v] not supported", method)
}

func logRequest(method string, request interface{}) {
	str := spew.Sdump(request)
	log.Printf("[%s] request: %s\n", method, str)
}

func NewDriver() *Driver {
	return &Driver{
		networks: make(map[string]*Network),
	}
}

func (t *Driver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	logRequest("GetCapabilities", nil)

	rootNs, err := netns.GetFromPid(1)
	if err != nil {
		return nil, err
	}
	t.rootNs = rootNs

	response := &network.CapabilitiesResponse{
		Scope:             network.LocalScope,
		ConnectivityScope: network.LocalScope,
	}
	log.Printf("[GetCapabilities] response: %+v\n", response)
	return response, nil
}

func (t *Driver) CreateNetwork(req *network.CreateNetworkRequest) error {
	logRequest("CreateNetwork", req)

	if len(req.IPv4Data) > 1 || len(req.IPv6Data) > 0 {
		return fmt.Errorf("Multiple ipv4 data or ipv6 data not supported")
	}

	options := req.Options["com.docker.network.generic"].(map[string]interface{})
	network, err := CreateNetwork(req.IPv4Data[0], options)
	if err != nil {
		return err
	}
	t.networks[req.NetworkID] = network

	return nil
}

func (t *Driver) AllocateNetwork(req *network.AllocateNetworkRequest) (*network.AllocateNetworkResponse, error) {
	logRequest("AllocateNetwork", req)
	return nil, notSupported("AllocateNetwork")
}

func (t *Driver) DeleteNetwork(req *network.DeleteNetworkRequest) error {
	logRequest("DeleteNetwork", req)
	return notSupported("DeleteNetwork")
}

func (t *Driver) FreeNetwork(req *network.FreeNetworkRequest) error {
	logRequest("FreeNetwork", req)
	return notSupported("FreeNetwork")
}

func (t *Driver) CreateEndpoint(req *network.CreateEndpointRequest) (*network.CreateEndpointResponse, error) {
	logRequest("CreateEndpoint", req)
	return nil, notSupported("CreateEndpoint")
}

func (t *Driver) DeleteEndpoint(req *network.DeleteEndpointRequest) error {
	logRequest("DeleteEndpoint", req)
	return notSupported("DeleteEndpoint")
}

func (t *Driver) EndpointInfo(req *network.InfoRequest) (*network.InfoResponse, error) {
	logRequest("EndpointInfo", req)
	return nil, notSupported("EndpointInfo")
}

func (t *Driver) Join(req *network.JoinRequest) (*network.JoinResponse, error) {
	logRequest("Join", req)
	return nil, notSupported("Join")
}

func (t *Driver) Leave(req *network.LeaveRequest) error {
	logRequest("Leave", req)
	return notSupported("Leave")
}

func (t *Driver) DiscoverNew(req *network.DiscoveryNotification) error {
	logRequest("DiscoverNew", req)
	return notSupported("DiscoverNew")
}

func (t *Driver) DiscoverDelete(req *network.DiscoveryNotification) error {
	logRequest("DiscoverDelete", req)
	return notSupported("DiscoverDelete")
}

func (t *Driver) ProgramExternalConnectivity(req *network.ProgramExternalConnectivityRequest) error {
	logRequest("ProgramExternalConnectivity", req)
	return notSupported("ProgramExternalConnectivity")
}

func (t *Driver) RevokeExternalConnectivity(req *network.RevokeExternalConnectivityRequest) error {
	logRequest("RevokeExternalConnectivity", req)
	return notSupported("RevokeExternalConnectivity")
}
