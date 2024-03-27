package model

import "net"

type IPAddr struct {
	ip net.IP
}

func NewIPAddr(ip string) *IPAddr {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return nil
	}
	return &IPAddr{ip: parsed}
}

func (s IPAddr) String() string {
	return s.ip.String()
}
