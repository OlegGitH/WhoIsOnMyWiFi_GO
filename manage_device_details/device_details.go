// device_details/device_details.go

package device_details

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// import necessary packages

func GetDeviceDetails() {
	// Logic to retrieve and store device details
}

func GetDeviceNameByMAC(mac string) (string, error) {
	// Execute the npcap's arp command to retrieve the ARP table
	cmd := exec.Command("arp", "-a")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse the output to find the device name associated with the MAC address
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, mac) {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[0], nil
			}
		}
	}

	// If the MAC address is not found in the ARP table
	return "", fmt.Errorf("device name for mac address %s not found", mac)
}

func FindDeviceName(routerIP, username, password, macAddress string) (string, error) {
	// Create HTTP client
	client := &http.Client{}

	// Prepare login request
	loginData := strings.NewReader("username=" + username + "&password=" + password)
	loginURL := "http://" + routerIP + "/login.cgi"
	loginReq, err := http.NewRequest("POST", loginURL, loginData)
	if err != nil {
		return "", err
	}

	// Send login request
	_, err = client.Do(loginReq)
	if err != nil {
		return "", err
	}

	// Prepare devices request
	devicesURL := "http://" + routerIP + "/devices.asp"
	devicesReq, err := http.NewRequest("GET", devicesURL, nil)
	if err != nil {
		return "", err
	}

	// Send devices request
	devicesResp, err := client.Do(devicesReq)
	if err != nil {
		return "", err
	}
	defer devicesResp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(devicesResp.Body)
	if err != nil {
		return "", err
	}

	// Find device name
	macClass := "mac_" + strings.Replace(macAddress, ":", "-", -1)
	startIdx := strings.Index(string(body), macClass)
	if startIdx == -1 {
		return "Device not found", nil
	}
	hostnameStart := strings.Index(string(body[startIdx:]), "<td class=\"hostname\">")
	hostnameEnd := strings.Index(string(body[startIdx+hostnameStart:]), "</td>")
	deviceName := string(body[startIdx+hostnameStart+len("<td class=\"hostname\">") : startIdx+hostnameStart+hostnameEnd])

	return deviceName, nil
}

func GetDeviceInformationPCAP() {
	// Find all available network interfaces
	// Find all available network interfaces
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Print available devices
	fmt.Println("Available devices:")
	for i, dev := range devices {
		fmt.Printf("[%d] Name: %s\n", i, dev.Name)
		fmt.Printf("    Description: %s\n", dev.Description)
		fmt.Println("------------------------------------------------------")
		// Open the network interface
		handle, err := pcap.OpenLive(dev.Name, 65536, true, pcap.BlockForever)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()

		// Get interface information
		netInterface, err := handle.Stats()
		if err != nil {
			log.Fatal(err)
		}
		// Convert the received packets count to a string
		packetsReceivedStr := strconv.Itoa(netInterface.PacketsReceived)

		// Print device information
		fmt.Println("netInterface.PacketsDropped", netInterface.PacketsDropped)
		fmt.Println("Interface Description:", netInterface.PacketsReceived)
		fmt.Println("Packets Received:", packetsReceivedStr)
	}

	// Choose an interface to capture packets from
	// For example, let's choose the first device
	if len(devices) > 0 {
		handle, err := pcap.OpenLive(devices[0].Name, 65536, true, pcap.BlockForever)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()

		// Set filter to capture only Wi-Fi management packets
		err = handle.SetBPFFilter("type mgt")
		if err != nil {
			log.Fatal(err)
		}

		// Start capturing packets
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		fmt.Println("Capturing Wi-Fi management packets...")
		for packet := range packetSource.Packets() {
			// Extract Wi-Fi management frame
			wifiLayer := packet.Layer(layers.LayerTypeDot11)
			if wifiLayer == nil {
				continue // Skip if not a Wi-Fi packet
			}
			wifiFrame, _ := wifiLayer.(*layers.Dot11)

			// Extract device information
			if wifiFrame.Type == layers.Dot11TypeMgmt {
				fmt.Printf("SSID: %s, BSSID: %s, Source: %s\n",
					wifiFrame.LayerContents(), wifiFrame.Address1, wifiFrame.Address2)
			}
		}
	}

	// Wait for a while to capture packets
	time.Sleep(30 * time.Second)
}
