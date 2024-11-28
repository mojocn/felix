package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func proxySettingOn(socks5Addr string) {
	networkService := "Wi-Fi"   // Change this to your active network service name
	proxyAddress := "127.0.0.1" // SOCKS5 proxy address
	proxyPort := "1080"         // SOCKS5 proxy port

	// Disable the SOCKS proxy first (optional cleanup)
	disableCmd := exec.Command("networksetup", "-setsocksfirewallproxystate", networkService, "off")
	if err := disableCmd.Run(); err != nil {
		fmt.Printf("Failed to disable SOCKS proxy: %v\n", err)
	}

	// Set the SOCKS proxy
	cmd := exec.Command("networksetup",
		"-setsocksfirewallproxy",
		networkService,
		proxyAddress,
		proxyPort)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to set SOCKS proxy: %v\n", err)
		return
	}

	// Enable the SOCKS proxy
	enableCmd := exec.Command("networksetup", "-setsocksfirewallproxystate", networkService, "on")
	if err := enableCmd.Run(); err != nil {
		fmt.Printf("Failed to enable SOCKS proxy: %v\n", err)
		return
	}

	fmt.Println("SOCKS5 proxy configured successfully!")
}

func proxySettingOff() (string, error) {
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// List all network services
	services := strings.Split(out.String(), "\n")
	for _, service := range services {
		if service != "" {
			// Check if the service is active by getting the status
			statusCmd := exec.Command("networksetup", "-getinfo", service)
			var statusOut bytes.Buffer
			statusCmd.Stdout = &statusOut
			if err := statusCmd.Run(); err != nil {
				return "", err
			}

			// Check if the service has a valid IP address (active network)
			if strings.Contains(statusOut.String(), "IP address") {
				return service, nil
			}
		}
	}
	return "", fmt.Errorf("no active network service found")
}
