package icmp

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
)

const (
	ProtocolNumber = 1
	ProtocolName   = "icmp"
)

func DialIP(localIP, remoteIP string) (*net.IPConn, error) {

	lAddr, err := net.ResolveIPAddr("ip", localIP)
	if err != nil {
		return nil, fmt.Errorf("local ip %s: %v", localIP, err)
	}

	rAddr, err := net.ResolveIPAddr("ip", remoteIP)
	if err != nil {
		return nil, fmt.Errorf("remote ip %s: %v", remoteIP, err)
	}
	return net.DialIP("ip4:"+ProtocolName, lAddr, rAddr)
}

func ReadPacketConn(c *icmp.PacketConn) (*icmp.Message, net.Addr, error) {

	b, addr, err := ReadPacketConnBytes(c)
	if err != nil {
		return nil, nil, fmt.Errorf("conn read from: %v", err)
	}

	msg, err := icmp.ParseMessage(ProtocolNumber, b)
	if err != nil {
		return nil, nil, fmt.Errorf("parse message: %v", err)
	}
	return msg, addr, nil
}

func ReadPacketConnBytes(c *icmp.PacketConn) ([]byte, net.Addr, error) {

	// icmp header 8 bytes, icmp body max. 576 bytes
	rb := make([]byte, 584)
	n, addr, err := c.ReadFrom(rb)
	if err != nil {
		return nil, nil, fmt.Errorf("conn read from: %v", err)
	}
	return rb[:n], addr, nil
}

func ReadIPConn(c *net.IPConn) (*ipv4.Header, *icmp.Message, error) {

	b, err := ReadIPConnBytes(c)
	if err != nil {
		return nil, nil, err
	}

	ipv4header, err := icmp.ParseIPv4Header(b)
	if err != nil {
		return nil, nil, fmt.Errorf("parse ipv4 header: %v", err)
	}

	msg, err := icmp.ParseMessage(ProtocolNumber, b[ipv4header.Len:])
	if err != nil {
		return nil, nil, fmt.Errorf("parse message: %v", err)
	}
	return ipv4header, msg, nil
}

func ReadIPConnBytes(c *net.IPConn) ([]byte, error) {

	// IP header 20 bytes, icmp header 8 bytes, icmp body max. 576 bytes
	rb := make([]byte, 604)
	n, err := c.Read(rb)
	if err != nil {
		return nil, fmt.Errorf("conn read: %v", err)
	}
	return rb[:n], nil
}
