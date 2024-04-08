// wifi_monitor/wifi_monitor.go

package wifi_monitoring

import (
	"fmt"
	"os/exec"
	"strings"
)

// import necessary packages

type WiFiConnection struct {
	Internet string
	Address  string
	Physical string
}

func MonitorWiFiConnections() error {
	// Execute the system command to retrieve the ARP table
	cmd := exec.Command("arp", "-a")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// Parse the output to extract MAC addresses and IP addresses
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && !strings.HasPrefix(fields[0], "Interface") {
			ip := fields[0]
			mac := fields[1]
			fmt.Printf("IP: %s, MAC: %s\n", ip, mac)
		}
	}

	return nil
}
