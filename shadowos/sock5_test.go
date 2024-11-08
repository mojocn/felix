package shadowos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net/http"
	"testing"
)

func httpSocks5Client() *http.Client {
	socks5Proxy := "127.0.0.1:2080"
	dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
	if err != nil {
		fmt.Printf("Failed to create SOCKS5 dialer: %v\n", err)
		return nil
	}

	// Create a custom transport that uses the SOCKS5 dialer
	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport: transport,
	}
	return client
}

func TestHttpOverSocks5(t *testing.T) {
	client := httpSocks5Client()

	// Generate a JSON body
	data := map[string]string{"deno": "felix"}
	jsonData, _ := json.Marshal(data)

	// URL to send the POST request to
	url := "https://httpbin.org/post"

	// Send the POST request
	response, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("failed to send POST request: %v", err)
		return
	}
	defer response.Body.Close()
	all, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failed to read response body: %v", err)
		return
	}
	t.Logf("Response: %s", all)

}

func TestAAAAA(t *testing.T) {
	list := []string{"a", "b", "c", "1"}
	aa := list[1:2]
	t.Log(aa)
}
