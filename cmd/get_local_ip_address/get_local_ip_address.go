package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	var localIP []net.IP

	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() != nil {
				localIP = append(localIP, ip)
			}
		}
	}

	for i, ip := range localIP {
		fmt.Printf("%d: %s\n", i, ip)
	}
}
