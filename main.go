package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func main() {
	// Get the local IP address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var localIP string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			localIP = ipnet.IP.String()
			break
		}
	}

	fmt.Println(localIP)
	// Run ARP command to get the list of devices
	cmd := exec.Command("arp", "-a")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Parse the ARP table and print the list of devices
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			fmt.Println(fields[0] + " : " + fields[1] + " : " + fields[2])
		}
	}
}
