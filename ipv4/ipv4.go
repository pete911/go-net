package ipv4

import (
	"fmt"
	"golang.org/x/net/ipv4"
	"net"
)

func Dial(protocol, localIP, remoteIP string) (*net.IPConn, error) {

	lAddr, err := net.ResolveIPAddr("ip", localIP)
	if err != nil {
		return nil, fmt.Errorf("local ip %s: %v", localIP, err)
	}

	rAddr, err := net.ResolveIPAddr("ip", remoteIP)
	if err != nil {
		return nil, fmt.Errorf("remote ip %s: %v", remoteIP, err)
	}

	network := "ip4:"
	if protocol != "" {
		network = network + protocol
	}

	return net.DialIP(network, lAddr, rAddr)
}

func ReadIPConn(c *net.IPConn, size int) (*ipv4.Header, []byte, error) {

	b, err := readIPConn(c, size)
	if err != nil {
		return nil, nil, err
	}
	return toIPv4Packet(b)
}

func toIPv4Packet(b []byte) (*ipv4.Header, []byte, error) {

	ipv4header, err := ipv4.ParseHeader(b)
	if err != nil {
		return nil, nil, fmt.Errorf("parse ipv4 header: %v", err)
	}
	return ipv4header, b[ipv4header.Len:], nil
}

func readIPConn(c *net.IPConn, size int) ([]byte, error) {

	rb := make([]byte, size)
	n, err := c.Read(rb)
	if err != nil {
		return nil, fmt.Errorf("conn read: %v", err)
	}
	return rb[:n], nil
}
