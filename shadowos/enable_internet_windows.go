package shadowos

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
)

func EnableInternetSetting(socks5Addr string) {
	// Open the registry key for proxy settings
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.SET_VALUE)
	if err != nil {
		log.Fatalf("Error opening registry key: %v", err)
	}
	defer key.Close()

	// Enable proxy and set proxy server to SOCKS5
	err = key.SetDWordValue("ProxyEnable", 1) // 1 to enable proxy
	if err != nil {
		log.Fatalf("Error enabling proxy: %v", err)
	}

	// Set the SOCKS5 proxy address (e.g., "socks=127.0.0.1:1080")
	err = key.SetStringValue("ProxyServer", "socks5://"+socks5Addr)
	if err != nil {
		log.Fatalf("Error setting proxy server: %v", err)
	}

	// Set the proxy override settings : *.cn;*.local
	//
	skipAddrs := "localhost;127.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*;192.168.*;<local>;*.cn"
	err = key.SetStringValue("ProxyOverride", skipAddrs)
	if err != nil {
		log.Fatalf("Error setting proxy override: %v", err)
	}
	// Optionally disable automatic proxy detection
	err = key.SetDWordValue("AutoDetect", 0) // 0 to disable
	if err != nil {
		log.Fatalf("Error disabling automatic proxy detection: %v", err)
	}

	fmt.Println("SOCKS5 proxy configuration applied successfully.")

	//cmd := exec.Command("netsh", "winhttp", "reset", "proxy")
	//err = cmd.Run()
	//if err != nil {
	//	slog.Error("Error resetting proxy settings: ", err)
	//}
	//log.Println("Network settings refreshed.")
}
func DisableInternetSetting() {
	log.Print("Disabling SOCKS5 proxy settings...")
	// Open the registry key for proxy settings
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.SET_VALUE)
	if err != nil {
		log.Fatalf("Error opening registry key: %v", err)
	}
	defer key.Close()

	// Enable proxy and set proxy server to SOCKS5
	err = key.SetDWordValue("ProxyEnable", 0) // 1 to enable proxy
	if err != nil {
		log.Fatalf("Error enabling proxy: %v", err)
	}

	//cmd := exec.Command("netsh", "winhttp", "reset", "proxy")
	//err = cmd.Run()
	//if err != nil {
	//	log.Fatalf("Error resetting proxy settings: %v", err)
	//}
	//log.Println("Network settings refreshed.")
}
