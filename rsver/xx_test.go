package rsver

import (
	"fmt"
	"net"
	"testing"
)

func TestName(t *testing.T) {

	// DNS server address
	dnsServer := "114.114.114.114:53"

	// Create a UDP connection
	conn, err := net.Dial("udp", dnsServer)
	if err != nil {
		fmt.Println("Error connecting to DNS server:", err)
		return
	}
	defer conn.Close()

	// Construct a DNS query for example.com
	query := []byte{
		0x12, 0x34, // Transaction ID
		0x01, 0x00, // Flags: standard query
		0x00, 0x01, // Questions: 1
		0x00, 0x00, // Answer RRs: 0
		0x00, 0x00, // Authority RRs: 0
		0x00, 0x00, // Additional RRs: 0
		0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, // Query: example
		0x03, 0x63, 0x6f, 0x6d, // Query: com
		0x00,       // Null terminator
		0x00, 0x01, // Type: A
		0x00, 0x01, // Class: IN
	}

	// Send the DNS query
	_, err = conn.Write(query)
	if err != nil {
		fmt.Println("Error sending DNS query:", err)
		return
	}

	// Set a read deadline

	// Read the response
	response := make([]byte, 512)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println("Error reading DNS response:", err)
		return
	}

	// Print the response
	fmt.Printf("Received %d bytes\n %q", n, response[:n])
}
