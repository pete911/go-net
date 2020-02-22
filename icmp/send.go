package icmp

import (
	"fmt"
	netipv4 "github.com/pete911/go-net/ipv4"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
)

func Send(localIP, remoteIP string, request *icmp.Message) (*ipv4.Header, *icmp.Message, error) {

	conn, err := netipv4.Dial(ProtocolName, localIP, remoteIP)
	if err != nil {
		return nil, nil, fmt.Errorf("ipv4 dial: %w", err)
	}
	defer conn.Close()

	wb, err := request.Marshal(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal request: %w", err)
	}

	if _, err := conn.Write(wb); err != nil {
		return nil, nil, fmt.Errorf("conn write: %v", err)
	}
	return ReadIPConn(conn)
}

func ReadIPConn(c *net.IPConn) (*ipv4.Header, *icmp.Message, error) {

	// IP header max 60 bytes, icmp header 8 bytes, icmp body max. 576 bytes
	ipv4header, data, err := netipv4.ReadIPConn(c, 644)
	if err != nil {
		return nil, nil, err
	}

	msg, err := icmp.ParseMessage(ProtocolNumber, data)
	if err != nil {
		return nil, nil, fmt.Errorf("parse message: %v", err)
	}
	return ipv4header, msg, nil
}
