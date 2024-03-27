package model

import (
	"net"
)

type Subnet struct {
	ipNet net.IPNet
}

func NewSubnet(subnet string) *Subnet {
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil
	}

	return &Subnet{ipNet: *ipnet}
}

func (s Subnet) String() string {
	return s.ipNet.String()
}

func (s Subnet) Equal(other Subnet) bool {
	return s.ipNet.String() == other.ipNet.String()
}

func (s Subnet) Contains(ip IPAddr) bool {
	return s.ipNet.Contains(ip.ip)
}
