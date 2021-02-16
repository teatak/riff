package common

import (
	"fmt"
	"net"
)

func AdviseRpc() (string, error) {
	addrs, err := GetPrivateIPv4()
	if err != nil {
		return "", err
	}
	if len(addrs) > 0 {
		return addrs[0].String(), nil
	}
	return "", err
}

func GetPrivateIPv4() ([]*net.IPNet, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	var ipNets []*net.IPNet
	for _, address := range addrs {
		// check the address type and if it is not a loopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if !isPrivate(ipnet.IP) {
					continue
				}
				ipNets = append(ipNets, ipnet)
			}
		}
	}
	return ipNets, nil
}

// privateBlocks contains non-forwardable address blocks which are used
// for private networks. RFC 6890 provides an overview of special
// address blocks.
var privateBlocks = []*net.IPNet{
	parseCIDR("10.0.0.0/8"),     // RFC 1918 IPv4 private network address
	parseCIDR("100.64.0.0/10"),  // RFC 6598 IPv4 shared address space
	parseCIDR("127.0.0.0/8"),    // RFC 1122 IPv4 loopback address
	parseCIDR("169.254.0.0/16"), // RFC 3927 IPv4 link local address
	parseCIDR("172.16.0.0/12"),  // RFC 1918 IPv4 private network address
	parseCIDR("192.0.0.0/24"),   // RFC 6890 IPv4 IANA address
	parseCIDR("192.0.2.0/24"),   // RFC 5737 IPv4 documentation address
	parseCIDR("192.168.0.0/16"), // RFC 1918 IPv4 private network address
	parseCIDR("::1/128"),        // RFC 1884 IPv6 loopback address
	parseCIDR("fe80::/10"),      // RFC 4291 IPv6 link local addresses
	parseCIDR("fc00::/7"),       // RFC 4193 IPv6 unique local addresses
	parseCIDR("fec0::/10"),      // RFC 1884 IPv6 site-local addresses
	parseCIDR("2001:db8::/32"),  // RFC 3849 IPv6 documentation address
}

func parseCIDR(s string) *net.IPNet {
	_, block, err := net.ParseCIDR(s)
	if err != nil {
		panic(fmt.Sprintf("Bad CIDR %s: %s", s, err))
	}
	return block
}

func isPrivate(ip net.IP) bool {
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}
