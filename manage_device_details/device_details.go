// device_details/device_details.go

package device_details

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// import necessary packages

func GetDeviceDetails() {
	// Logic to retrieve and store device details
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
