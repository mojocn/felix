package socks5ws

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"sort"
	"sync"
	"time"
)

type RelayUdpDirect struct {
	s5       net.Conn
	relayUdp *net.UDPConn

	// Reassembly queue for fragmented UDP packets.
	mu            sync.Mutex
	fragments     map[string][]*udpPacket // Map of DstAddr to fragments
	highestFrag   map[string]byte         // Track highest FRAG value for each DstAddr
	timers        map[string]*time.Timer  // Map of DstAddr to reassembly timer
	timerDuration time.Duration           // Timer duration
}

func (ud *RelayUdpDirect) addFragment(clientAddr *net.UDPAddr, frag *udpPacket) {
	ud.mu.Lock()
	defer ud.mu.Unlock()
	clientDstAddr := ud.clientDstAddrAsID(clientAddr, frag.dstAddr())
	// Initialize fragment queue and timer if not already present.
	if _, exists := ud.fragments[clientDstAddr]; !exists {
		ud.fragments[clientDstAddr] = []*udpPacket{}
		ud.highestFrag[clientDstAddr] = socks5UdpFragNotSupported
		ud.startTimer(clientDstAddr)
	}

	// Update highest FRAG value.
	if frag.Frag > ud.highestFrag[clientDstAddr] {
		ud.highestFrag[clientDstAddr] = frag.Frag
	}

	// Add fragment to the queue.
	ud.fragments[clientDstAddr] = append(ud.fragments[clientDstAddr], frag)

	// Check if this is the final fragment (end-of-fragment sequence).
	if frag.Frag == socks5UdpFragEnd || frag.Frag == socks5UdpFragNotSupported { // High-order bit indicates end of sequence.
		ud.assembleThenPipeUdp(clientAddr, frag.dstAddr())
	}
}

func (ud *RelayUdpDirect) startTimer(ClientDstAddr string) {
	if timer, exists := ud.timers[ClientDstAddr]; exists {
		timer.Stop()
	}
	ud.timers[ClientDstAddr] = time.AfterFunc(ud.timerDuration, func() {
		ud.mu.Lock()
		defer ud.mu.Unlock()
		delete(ud.fragments, ClientDstAddr)
		delete(ud.highestFrag, ClientDstAddr)
		delete(ud.timers, ClientDstAddr)
	})
}
func (ud *RelayUdpDirect) clientDstAddrAsID(clientAddr *net.UDPAddr, dstAddr string) string {
	return fmt.Sprintf("%s/%s", clientAddr, dstAddr)
}
func (ud *RelayUdpDirect) assembleThenPipeUdp(clientAddr *net.UDPAddr, dstAddr string) {
	var data []byte
	clientDstAddr := ud.clientDstAddrAsID(clientAddr, dstAddr)
	fragments := ud.fragments[clientDstAddr]
	// Sort fragments by FRAG value.
	sort.Slice(fragments, func(i, j int) bool {
		return fragments[i].Frag < fragments[j].Frag
	})
	for _, frag := range fragments {
		data = append(data, frag.Data...)
	}
	comboPacket := fragments[0]
	comboPacket.Data = data

	// Clean up after successful reassembly.
	delete(ud.fragments, clientDstAddr)
	delete(ud.highestFrag, clientDstAddr)
	if timer, exists := ud.timers[clientDstAddr]; exists {
		timer.Stop()
		delete(ud.timers, clientDstAddr)
	}
	ud.segmentPipe(comboPacket, clientAddr)
}

func (ud *RelayUdpDirect) StartPipe() {
	buf := make([]byte, udpMTU)
	for {
		n, clientAddr, err := ud.relayUdp.ReadFromUDP(buf)
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			slog.Error("Error reading UDP data", "err", err.Error())
			continue
		}
		packet, err := parseUDPData(buf[:n])
		if err != nil {
			log.Println("Error parsing UDP data", err)
			continue
		}
		ud.addFragment(clientAddr, packet)
	}
}

func (ud *RelayUdpDirect) segmentPipe(comboPacket *udpPacket, clientAddr *net.UDPAddr) {
	resp, err := forwardUDPData(comboPacket)
	if err != nil {
		slog.Error("Error forwarding UDP data", "err", err.Error())
		return
	}
	header := comboPacket.ResponseData(resp)
	_, err = ud.relayUdp.WriteToUDP(header, clientAddr)
	if err != nil {
		slog.Error("Error sending UDP response", "err", err.Error())
	}
}

func NewRelayUdpDirect(s5 net.Conn) (*RelayUdpDirect, error) {
	udpAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to bind UDP socket: %w", err)
	}
	ud := &RelayUdpDirect{
		s5:            s5,
		relayUdp:      udpConn,
		mu:            sync.Mutex{},
		fragments:     make(map[string][]*udpPacket),
		highestFrag:   make(map[string]byte),
		timers:        make(map[string]*time.Timer),
		timerDuration: time.Second * 60,
	}

	boundAddr := udpConn.LocalAddr().(*net.UDPAddr)
	response := []byte{
		socks5Version, socks5ReplyOkay, socks5ReplyReserved, socks5AtypeIPv4,
		boundAddr.IP[0], boundAddr.IP[1], boundAddr.IP[2], boundAddr.IP[3],
		byte(boundAddr.Port >> 8), byte(boundAddr.Port & 0xFF),
	}
	_, err = ud.s5.Write(response)
	if err != nil {
		return nil, fmt.Errorf("failed to send response to client: %w", err)
	}

	return ud, nil
}

func forwardUDPData(udpPacket *udpPacket) ([]byte, error) {
	conn, err := net.DialUDP("udp", nil, udpPacket.addr())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(udpPacket.Data)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, udpMTU)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

type udpPacket struct {
	RSV   [2]byte // 保留字段
	Frag  byte    // 分片字段
	AType byte    // dst 地址类型
	Addr  []byte  // dst 地址
	Port  []byte  // dst 端口
	Data  []byte  // 数据
}

func (p udpPacket) ResponseData(payload []byte) []byte {
	header := []byte{p.RSV[0], p.RSV[1], 0, p.AType}
	header = append(header, p.Addr...)
	header = append(header, p.Port...)
	return append(header, payload...)
}

func parseUDPData(data []byte) (*udpPacket, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid UDP packet")
	}
	// 解析头部
	var packet = udpPacket{
		RSV:   [2]byte{data[0], data[1]},
		Frag:  data[2],
		AType: data[3],
	}
	switch packet.AType {
	case socks5AtypeIPv4:
		if len(data) < 10 {
			return nil, fmt.Errorf("invalid IPv4 UDP packet")
		}
		packet.Addr = data[4 : 4+net.IPv4len]
		packet.Port = data[4+net.IPv4len : 4+net.IPv4len+2]
		packet.Data = data[4+net.IPv4len+2:]
	case socks5AtypeIPv6:
		if len(data) < 22 {
			return nil, fmt.Errorf("invalid IPv6 UDP packet")
		}
		packet.Addr = data[4 : 4+net.IPv6len]
		packet.Port = data[4+net.IPv6len : 4+net.IPv6len+2]
		packet.Data = data[4+net.IPv6len+2:]
	case socks5AtypeDomain:
		if len(data) < 7 {
			return nil, fmt.Errorf("invalid domain UDP packet")
		}
		addrLen := int(data[4])
		packet.Addr = data[5 : 5+addrLen]
		packet.Port = data[5+addrLen : 5+addrLen+2]
		packet.Data = data[5+addrLen+2:]
	default:
		return nil, fmt.Errorf("unsupported address type: %d", packet.AType)
	}
	// 返回目标地址和数据
	return &packet, nil
}

func (p udpPacket) ip() net.IP {
	switch p.AType {
	case socks5AtypeIPv4, socks5AtypeIPv6:
		return p.Addr
	case socks5AtypeDomain:
		ips, err := net.LookupIP(string(p.Addr))
		if err != nil {
			slog.Error("failed to resolve domain", "err", err.Error())
			return net.IPv4zero
		}
		if len(ips) == 0 {
			return net.IPv4zero
		}
		return ips[0]
	default:
		return net.IPv4zero
	}
}

func (p udpPacket) port() int {
	return int(p.Port[0])<<8 + int(p.Port[1])
}
func (p udpPacket) addr() *net.UDPAddr {
	return &net.UDPAddr{IP: p.ip(), Port: p.port()}
}
func (p udpPacket) dstAddr() string {
	return fmt.Sprintf("%s:%d", p.ip(), p.port())
}

func (ud *RelayUdpDirect) Close() {
	//s5 has already been closed in outside
	if ud.relayUdp != nil {
		err := ud.relayUdp.Close()
		if err != nil {
			log.Println("close udp conn failed: ", err)
		}
	}

}
