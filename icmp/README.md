# ICMP

## ICMP Header

```
+-------------------------------+-------------------------------+-------------------------------+-------------------------------+
|               0               |               1               |               2               |               3               |
+-------------------------------+-------------------------------+-------------------------------+-------------------------------+
| 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10| 11| 12| 13| 14| 15| 16| 17| 18| 19| 20| 21| 22| 23| 24| 25| 26| 27| 28| 29| 30| 31|
+-------------------------------+-------------------------------+---------------------------------------------------------------+
|           Type                |            Code               |                            Checksum                           |
+-------------------------------+-------------------------------+---------------------------------------------------------------+
|                                                         Rest of Header                                                        |
+-------------------------------------------------------------------------------------------------------------------------------+
```

Header is 8 bytes (1 byte type, 1 byte code, 2 bytes checksum and 4 bytes rest of header). Data is typically 56 bytes
(whole ICMP packet including header 64 bytes) and maximum 576 bytes (whole ICMP packet including header 584 bytes).

Because IP packet can be 65,535 bytes in size, ICMP body can be set to 65,507 (65,535 - 20 (IP header) - 8 (ICMP header)).

## ICMP Ping request example

```go
import (
	neticmp "github.com/pete911/go-net/icmp"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"os"
)

func IcmpPing() {

	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("hello ..."),
		},
	}

	ipv4Header, icmpResponse, err := neticmp.Send("127.0.0.1", "127.0.0.1", msg)
	if err != nil {
		log.Fatalf("icmp ping: %v", err)
	}
	b, err := icmpResponse.Body.Marshal(neticmp.ProtocolNumber)
	if err != nil {
		log.Printf("icmp body marshal: %v", err)
		return
	}

	log.Printf("IPv4 header: %s", ipv4Header)
	log.Printf("icmp: %+v", icmpResponse)
	log.Printf("icmp type: %d code: %d checksum: %d", icmpResponse.Type, icmpResponse.Code, icmpResponse.Checksum)
	log.Printf("icmp rest of header: %v", b[:4])
	log.Printf("icmp body: %s", string(b[4:])) // first 4 bytes is 'rest of header'
}
```

## ICMP listen example

```go
import (
	neticmp "github.com/pete911/go-net/icmp"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
)

func IcmpListen() {

	conns, err := neticmp.Listen("127.0.0.1")
	if err != nil {
		log.Fatalf("icmp listen: %v", err)
	}

	for conn := range conns {
		icmpHandler(conn)
	}
}

func icmpHandler(conn *icmp.PacketConn) {

	msg, addr, err := neticmp.ReadPacketConn(conn)
	if err != nil {
		log.Fatal(err)
	}

	if msg.Type != ipv4.ICMPTypeEcho {
		log.Printf("icmp ping from %s:", addr)
		log.Printf("%+v", *msg)
		return
	}

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("hello ..."),
		},
	}

	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := conn.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(addr.String())}); err != nil {
		log.Fatalf("icmp write to: %v", err)
	}
}
```