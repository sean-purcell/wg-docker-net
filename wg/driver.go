package wg

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/go-plugins-helpers/network"
)

type Driver struct {
	network.Driver
}

func notSupported(method string) error {
	return fmt.Errorf("[%v] not supported", method)
}
func logRequest(method string, request interface{}) {
	str := spew.Sdump(request)
	log.Printf("[%s] request: %s\n", method, str)
}

func (t *Driver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	logRequest("GetCapabilities", nil)
	response := &network.CapabilitiesResponse{
		Scope:             network.LocalScope,
		ConnectivityScope: network.LocalScope,
	}
	log.Printf("[GetCapabilities] response: %+v\n", response)
	return response, nil
}

func (t *Driver) CreateNetwork(req *network.CreateNetworkRequest) error {
	logRequest("CreateNetwork", req)
	return notSupported("CreateNetwork")
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
