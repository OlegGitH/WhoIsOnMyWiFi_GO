package main

import (
	"flag"
	"fmt"
	"os"
	device_details "whoisonmywifi/manage_device_details"
	"whoisonmywifi/wifi_monitoring"
)

func main() {

	command := flag.NewFlagSet("getpcapinfo", flag.ExitOnError)
	// Define flags for other commands if needed...

	// Check if there are any command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: whois <command>")
		os.Exit(1)
	}

	// Parse the command-line arguments
	switch os.Args[1] {
	case "getpcapinfo":
		command.Parse(os.Args[2:])
		device_details.GetDeviceInformationPCAP()
		// Add cases for other commands if needed...
	case "checkwifi":
		command.Parse(os.Args[2:])
		// Function to to get WiFi Connection data
		err, w := wifi_monitoring.MonitorWiFiConnections()

		if err != nil {
			return
		}

		for _, t := range *w {
			fmt.Printf("IP: %s, MAC: %s, Type: %s\n", t.IP, t.Mac, t.Type)
			// Get Device name from Router
			result, _ := device_details.GetDeviceNameByMAC(t.Mac)

			fmt.Print(result)
		}
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}

}
