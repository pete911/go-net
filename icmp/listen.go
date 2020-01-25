package icmp

import (
	"golang.org/x/net/icmp"
	"log"
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
