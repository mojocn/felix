package util

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net"
)

type SchemaVLESS struct {
	userID      uuid.UUID
	DstProtocol string //tcp or udp
	dstHost     string
	dstHostType string //ipv6 or ipv4,domain
	dstPort     uint16
	Version     byte
	playload    []byte
}

func (h SchemaVLESS) UUID() string {
	return h.userID.String()
}

func (h SchemaVLESS) DataUdp() []byte {
	allData := make([]byte, 0)
	chunk := h.playload
	for index := 0; index < len(chunk); {
		if index+2 > len(chunk) {
			fmt.Println("Incomplete length buffer")
			return nil
		}
		lengthBuffer := chunk[index : index+2]
		udpPacketLength := binary.BigEndian.Uint16(lengthBuffer)
		if index+2+int(udpPacketLength) > len(chunk) {
			fmt.Println("Incomplete UDP packet")
			return nil
		}
		udpData := chunk[index+2 : index+2+int(udpPacketLength)]
		index = index + 2 + int(udpPacketLength)
		allData = append(allData, udpData...)
	}
	return allData
}
func (h SchemaVLESS) DataTcp() []byte {
	return h.playload
}

func (h SchemaVLESS) AddrUdp() *net.UDPAddr {
	return &net.UDPAddr{IP: h.HostIP(), Port: int(h.dstPort)}
}
func (h SchemaVLESS) HostIP() net.IP {
	ip := net.ParseIP(h.dstHost)
	if ip == nil {
		ips, err := net.LookupIP(h.dstHost)
		if err != nil {
			h.Logger().Error("failed to resolve domain", "err", err.Error())
			return net.IPv4zero
		}
		if len(ips) == 0 {
			return net.IPv4zero
		}
		return ips[0]
	}
	return ip
}

func (h SchemaVLESS) HostPort() string {
	return net.JoinHostPort(h.dstHost, fmt.Sprintf("%d", h.dstPort))
}
func (h SchemaVLESS) Logger() *slog.Logger {
	return slog.With("userID", h.userID.String(), "network", h.DstProtocol, "addr", h.HostPort())
}

// VlessParse https://xtls.github.io/development/protocols/vless.html
func VlessParse(buf []byte) (*SchemaVLESS, error) {
	payload := &SchemaVLESS{
		userID:      uuid.Nil,
		DstProtocol: "",
		dstHost:     "",
		dstPort:     0,
		Version:     0,
		playload:    nil,
	}

	if len(buf) < 24 {
		return payload, errors.New("invalid payload length")
	}

	payload.Version = buf[0]
	payload.userID = uuid.Must(uuid.FromBytes(buf[1:17]))
	extraInfoProtoBufLen := buf[17]

	command := buf[18+extraInfoProtoBufLen]
	switch command {
	case 1:
		payload.DstProtocol = "tcp"
	case 2:
		payload.DstProtocol = "udp"
	default:
		return payload, fmt.Errorf("command %d is not supported, command 01-tcp, 02-udp, 03-mux", command)
	}

	portIndex := 18 + extraInfoProtoBufLen + 1
	payload.dstPort = binary.BigEndian.Uint16(buf[portIndex : portIndex+2])

	addressIndex := portIndex + 2
	addressType := buf[addressIndex]
	addressValueIndex := addressIndex + 1

	switch addressType {
	case 1: // IPv4
		if len(buf) < int(addressValueIndex+net.IPv4len) {
			return nil, fmt.Errorf("invalid IPv4 address length")
		}
		payload.dstHost = net.IP(buf[addressValueIndex : addressValueIndex+net.IPv4len]).String()
		payload.playload = buf[addressValueIndex+net.IPv4len:]
		payload.dstHostType = "ipv4"
	case 2: // domain
		addressLength := buf[addressValueIndex]
		addressValueIndex++
		if len(buf) < int(addressValueIndex)+int(addressLength) {
			return nil, fmt.Errorf("invalid domain address length")
		}
		payload.dstHost = string(buf[addressValueIndex : int(addressValueIndex)+int(addressLength)])
		payload.playload = buf[int(addressValueIndex)+int(addressLength):]
		payload.dstHostType = "domain"

	case 3: // IPv6
		if len(buf) < int(addressValueIndex+net.IPv6len) {
			return nil, fmt.Errorf("invalid IPv6 address length")
		}
		payload.dstHost = net.IP(buf[addressValueIndex : addressValueIndex+net.IPv6len]).String()
		payload.playload = buf[addressValueIndex+net.IPv6len:]
		payload.dstHostType = "ipv6"
	default:
		return nil, fmt.Errorf("addressType %d is not supported", addressType)
	}

	return payload, nil
}
