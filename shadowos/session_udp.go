package shadowos

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
)

type SessionUdp struct {
	SessionTcp
	udpConn *net.UDPConn
}

func (st *SessionUdp) breakGfwSvr(cfg *ProxyCfg) error {
	udpAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	udpConn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		return fmt.Errorf("failed to bind UDP socket: %w", err)
	}
	st.udpConn = udpConn

	boundAddr := udpConn.LocalAddr().(*net.UDPAddr)
	response := []byte{
		socks5Version, socks5ReplySuccess, socks5ReplyReserved, socks5AtypeIPv4,
		boundAddr.IP[0], boundAddr.IP[1], boundAddr.IP[2], boundAddr.IP[3],
		byte(boundAddr.Port >> 8), byte(boundAddr.Port & 0xFF),
	}

	wsAddr := cfg.WsUrl
	ws, _, err := websocket.DefaultDialer.Dial(wsAddr, cfg.WsHeader)
	if err != nil {
		return fmt.Errorf("failed to connect to remote proxy server: %s ,error:%v", wsAddr, err)
	}
	// Send success response
	_, err = st.s5.Write(response)
	if err != nil {
		log.Printf("failed to send response to client: %v", err)
	}
	ws.SetCloseHandler(func(code int, text string) error {
		log.Println("ws closed", code, text)
		st.wsExitCh <- struct{}{}
		return nil
	})
	st.ws = ws
	return nil
}

func (st *SessionUdp) Close() {
	//s5 has already been closed in outside
	if st.udpConn != nil {
		err := st.udpConn.Close()
		if err != nil {
			log.Println("close udp conn failed: ", err)
		}
	}
	if ws := st.ws; ws != nil {
		err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("send websocket close message failed: ", err)
		}
		err = ws.Close()
		if err != nil {
			log.Println("close websocket conn failed: ", err)
		}
	}
}

func (st *SessionUdp) pipe(ctx context.Context, uid [16]byte) {
	//udp is not working
	exitCh := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				exitCh <- struct{}{}
				return
			case <-st.wsExitCh:
				exitCh <- struct{}{}
				return
			default:
				packet := make([]byte, 65535)
				n, clientAddr, err := st.udpConn.ReadFromUDP(packet)
				if err != nil {
					log.Printf("UDP read error: %v", err)
					return
				}
				if n > 0 {
					st.handleUdp53Packet(st.udpConn, clientAddr, packet[:n], uid)
				}
			}
		}
	}()
	buf := make([]byte, 1)
	n, err := st.s5.Read(buf) // Block until client closes the connection
	if err != nil && errors.Is(err, io.EOF) {
		log.Printf("Failed to read from client: %v", err)
	} else {
		log.Println("s5 Client closed connection", n)
	}
	exitCh <- struct{}{}
}

func (st *SessionUdp) handleUdp53Packet(conn *net.UDPConn, clientAddr *net.UDPAddr, udpPacket []byte, uid [16]byte) {
	// Parse UDP udpPacket
	frag := udpPacket[2]
	atyp := udpPacket[3]
	if frag != 0x00 {
		log.Printf("Fragmented UDP packets are not supported")
		return
	}
	st.req.socks5Atyp = atyp
	if atyp != socks5AtypeIPv4 && atyp != socks5AtypeIPv6 && atyp != socks5AtypeDomain {
		log.Printf("Unsupported address type: %v", atyp)
		return
	}

	dstPortIndex := 4
	if atyp == socks5AtypeIPv4 {
		st.req.dstAddr = udpPacket[4 : 4+net.IPv4len]
		dstPortIndex += net.IPv4len
	} else if atyp == socks5AtypeIPv6 {
		st.req.dstAddr = udpPacket[4 : 4+net.IPv6len]
		dstPortIndex += net.IPv6len
	} else if atyp == socks5AtypeDomain {
		addrLen := int(udpPacket[4])
		st.req.dstAddr = udpPacket[5 : 5+addrLen]
		dstPortIndex += addrLen + 1
	} else {
		log.Printf("Unsupported address type: %v", atyp)
		return
	}
	st.req.dstPort = udpPacket[dstPortIndex : dstPortIndex+2]

	header, err := st.req.vlessHeaderUdp(uid)
	if err != nil {
		log.Println("failed to generate vless header")
		return
	}
	payload := udpPacket[dstPortIndex+2:]
	payloadN := len(payload)
	//payloadN to 2 bytes
	payload = append([]byte{byte(payloadN >> 8), byte(payloadN & 0xFF)}, payload...)

	data := append(header, payload...)

	err = st.ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Printf("failed to send UDP udpPacket to remote server: %v", err)
		return
	}
	mt, p, err := st.ws.ReadMessage()
	if err != nil {
		log.Printf("failed to read response from remote server: %v  %v", mt, err)
		return
	}
	if len(p) > 2 {
		fromIdx := p[1] + 2
		res := p[fromIdx:]
		socks5UdpHeader := []byte{0x00, 0x00, 0x00}
		if st.req.socks5Atyp == socks5AtypeIPv4 {
			socks5UdpHeader = append(socks5UdpHeader, 0x01)
			socks5UdpHeader = append(socks5UdpHeader, st.req.dstAddr...)
		} else if st.req.socks5Atyp == socks5AtypeIPv6 {
			socks5UdpHeader = append(socks5UdpHeader, 0x04)
			socks5UdpHeader = append(socks5UdpHeader, st.req.dstAddr...)
		} else if st.req.socks5Atyp == socks5AtypeDomain {
			socks5UdpHeader = append(socks5UdpHeader, 0x03)
			socks5UdpHeader = append(socks5UdpHeader, byte(len(st.req.dstAddr)))
			socks5UdpHeader = append(socks5UdpHeader, st.req.dstAddr...)
		} else {
			log.Println("Unsupported address type")
			return
		}
		socks5UdpHeader = append(socks5UdpHeader, st.req.dstPort...)

		res = append(socks5UdpHeader, res...)
		if _, err := conn.WriteToUDP(res, clientAddr); err != nil {
			log.Printf("Failed to send response to %v: %v", clientAddr, err)
		}
	}

	// Forward data (implement logic for forwarding here)

	// Example: Echo back to client

}
