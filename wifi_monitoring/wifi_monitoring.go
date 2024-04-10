// wifi_monitor/wifi_monitor.go

package wifi_monitoring

import (
	"os/exec"
	"strings"
)

// import necessary packages

type WiFiConnections struct {
	IP   string
	Mac  string
	Type string
}

func MonitorWiFiConnections() (*error, *[]WiFiConnections) {
	// Execute the system command to retrieve the ARP table
	cmd := exec.Command("arp", "-a")
	output, err := cmd.Output()
	if err != nil {
		return &err, nil
	}

	wifiConnectionTable := []WiFiConnections{}

	// Parse the output to extract MAC addresses and IP addresses
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && !strings.HasPrefix(fields[0], "Interface") && !strings.Contains(fields[0], "Internet") {
			wifiConnection := WiFiConnections{
				IP:   fields[0],
				Mac:  fields[1],
				Type: fields[2],
			}
			wifiConnectionTable = append(wifiConnectionTable, wifiConnection)

		}
	}

	return nil, &wifiConnectionTable
}
