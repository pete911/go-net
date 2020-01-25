# ICMP

## ICMP Header

```
+-------+-------+-------+-------+
|   0   |   1   |   2   |   3   |
+-------+-------+-------+-------+
| Type  | Code  | Checksum      |
+-------+-------+---------------+
| Rest of Header                |
+-------------------------------+
```

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