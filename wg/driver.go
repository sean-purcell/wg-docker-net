package wg

import (
	"fmt"
	"log"

	"github.com/docker/go-plugins-helpers/network"
)

type Driver struct {
	network.Driver
}

func (t *Driver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	return nil, fmt.Errorf("not supported")
}
