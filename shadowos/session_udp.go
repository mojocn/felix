package shadowos

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
	"log"
	"net"
)

type SessionUdp struct {
	req *Socks5Request
	s5  net.Conn
	ws  *websocket.Conn
}

func handleUDPPacket(conn *net.UDPConn, clientAddr *net.UDPAddr, packet []byte) {
	// Parse UDP packet
	headerLen := 3 + int(packet[4]) // Assuming SOCKS5 UDP header
	data := packet[headerLen:]

	frag := packet[2]
	atyp := packet[3]
	if frag != 0x00 {
		log.Printf("Fragmented UDP packets are not supported")
		return
	}
	if atyp != socks5AtypeIPv4 && atyp != socks5AtypeIPv6 && atyp != socks5AtypeDomain {
		log.Printf("Unsupported address type: %v", atyp)
		return
	}
	dstAddr := packet[4:]
	dstPortIndex := 4
	if atyp == socks5AtypeIPv4 {
		dstAddr = packet[4 : 4+net.IPv4len]
		dstPortIndex += net.IPv4len
	} else if atyp == socks5AtypeIPv6 {
		dstAddr = packet[4 : 4+net.IPv6len]
		dstPortIndex += net.IPv6len
	} else if atyp == socks5AtypeDomain {
		addrLen := int(packet[4])
		dstAddr = packet[5 : 5+addrLen]
		dstPortIndex += addrLen + 1
	} else {
		log.Printf("Unsupported address type: %v", atyp)
		return
	}
	dstPort := binary.BigEndian.Uint16(packet[dstPortIndex : dstPortIndex+2])

	payload := packet[dstPortIndex+2:]
	log.Printf("Received UDP packet for %v:%v from %v %v", dstAddr, dstPort, clientAddr, payload)

	targetAddr := packet[3:]
	log.Printf("Received UDP packet for %v from %v", targetAddr, clientAddr)

	// Forward data (implement logic for forwarding here)

	// Example: Echo back to client
	if _, err := conn.WriteToUDP(data, clientAddr); err != nil {
		log.Printf("Failed to send response to %v: %v", clientAddr, err)
	}
}
