package main

import (
	"log"
	"net"
)

// Function to perform DNS lookup if the host is not an IP address
func resolveHost(host string) string {
	ip := net.ParseIP(host)
	if ip == nil {
		// Host is not an IP address, perform DNS lookup
		addrs, err := net.LookupHost(host)
		if err != nil {
			log.Fatalf("Error resolving host: %s", err)
		}
		// Return the first resolved IP address
		return addrs[0]
	}
	// Host is already an IP address
	return host
}
