package main

import (
	"fmt"
	device_details "whoisonmywifi/manage_device_details"
	"whoisonmywifi/wifi_monitoring"
)

func main() {

	// Function to to get WiFi Connection data
	err, w := wifi_monitoring.MonitorWiFiConnections()

	if err != nil {
		return
	}

	for _, t := range *w {
		fmt.Printf("IP: %s, MAC: %s, Type: %s\n", t.IP, t.Mac, t.Type)
		// Get Device name from Router
		result, _ := device_details.GetDeviceNameByMAC(t.Mac)
		device_details.GetDeviceInformationPCAP()
		fmt.Print(result)
	}

}
