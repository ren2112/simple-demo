package utils

import (
	"fmt"
	"net"
)

func GetLocalIPv4() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		if i.Name == "WLAN" {
			addrs, err := i.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("WLAN IPv4 address not found")
}
