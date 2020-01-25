package icmp

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func Send(localIP, remoteIP string, request *icmp.Message) (*ipv4.Header, *icmp.Message, error) {

	conn, err := DialIP(localIP, remoteIP)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	wb, err := request.Marshal(nil)
	if err != nil {
		return nil, nil, err
	}

	if _, err := conn.Write(wb); err != nil {
		return nil, nil, fmt.Errorf("conn write: %v", err)
	}
	return ReadIPConn(conn)
}
