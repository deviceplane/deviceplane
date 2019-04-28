package info

import (
	"net"
)

func getIPAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, g := range addrs {
		if ipnet, ok := g.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ipnet4 := ipnet.IP.To4()
			if ipnet4 != nil {
				return ipnet4.String(), nil
			}
		}
	}

	return "", nil
}
