package icmp

import (
	"fmt"
	"golang.org/x/net/icmp"
	"log"
	"net"
	"os"
)

func Listen(ip string) (chan *icmp.PacketConn, error) {

	conns := make(chan *icmp.PacketConn)
	go func() {
		for {
			conn, err := icmp.ListenPacket("ip4:"+ProtocolName, ip)
			if err != nil {
				log.Printf("listen ip: %v", err)
				if err := conn.Close(); err != nil {
					log.Printf("conn close: %v", err)
				}
				close(conns)
				os.Exit(1)
			}
			conns <- conn
		}
	}()
	return conns, nil
}

func Read(c *icmp.PacketConn) (*icmp.Message, net.Addr, error) {

	// icmp header 8 bytes, icmp body max. 576 bytes
	b, addr, err := readFromPacketConn(c, 584)
	if err != nil {
		return nil, nil, fmt.Errorf("conn read from: %v", err)
	}

	msg, err := icmp.ParseMessage(ProtocolNumber, b)
	if err != nil {
		return nil, nil, fmt.Errorf("parse message: %v", err)
	}
	return msg, addr, nil
}

func readFromPacketConn(c *icmp.PacketConn, size int) ([]byte, net.Addr, error) {

	rb := make([]byte, size)
	n, addr, err := c.ReadFrom(rb)
	if err != nil {
		return nil, nil, fmt.Errorf("conn read from: %v", err)
	}
	return rb[:n], addr, nil
}
