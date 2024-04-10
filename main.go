package main

import (
	"fmt"
	device_details "whoisonmywifi/manage_device_details"
	"whoisonmywifi/wifi_monitoring"
)

var (
	RouterIP = "192.168.100.1"
	UserName = "root"
	Password = "admin"
)

func main() {

	// Function to to get WiFi Connection data
	err, w := wifi_monitoring.MonitorWiFiConnections()

	if err != nil {
		return
	}

	for _, t := range *w {
		fmt.Printf("IP: %s, MAC: %s, Type: %s\n", t.IP, t.Mac, t.Type)
		result, _ := device_details.FindDeviceName(RouterIP, UserName, Password, t.Mac)
		fmt.Printf(result)
	}

}
