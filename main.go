package main

import (
	"encoding/binary"
	"fmt"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"runtime"
)

// ping example with spoofed mac and IP

const etherTypeIPv4 uint16 = 0x0800

var (
	// update accordingly
	interfaceName = "en0"
	srcMAC = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	dstMAC = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	srcIP  = net.ParseIP("192.168.1.1")
	dstIP  = net.ParseIP("192.168.1.2")
)

func main() {

	ifi, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Fatalf("interface by name %s: %v", interfaceName, err)
	}

	c, err := raw.ListenPacket(ifi, etherTypeIPv4, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pingPacket, err := getICMPPingPacket("hello ...")
	if err != nil {
		log.Fatalf("get icmp ping packet: %v", err)
	}
	log.Printf("ping packet len: %d", len(pingPacket))

	ipPacket, err := getIPPacket(srcIP, dstIP, 1, pingPacket)
	if err != nil {
		log.Fatalf("get ip packet: %v", err)
	}
	log.Printf("ip packet len: %d", len(ipPacket))

	sendMessage(c, srcMAC, dstMAC, etherTypeIPv4, ipPacket)
}

func getICMPPingPacket(data string) ([]byte, error) {

	ping := &icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte(data),
		},
	}
	return ping.Marshal(nil)
}

func getIPPacket(src, dst net.IP, protocol int, msg []byte) ([]byte, error) {

	iph := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TOS:      0,
		TotalLen: ipv4.HeaderLen + len(msg),
		TTL:      64,
		Protocol: protocol,
		Dst:      dst,
		Src:      src,
	}

	ip, err := iph.Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshal ip request: %w", err)
	}

	// this is currently broken in golang, need to set total len correctly
	if runtime.GOOS == "darwin" {
		binary.BigEndian.PutUint16(ip[2:4], uint16(iph.TotalLen))
	}
	return append(ip, msg...), nil
}

func sendMessage(c net.PacketConn, src net.HardwareAddr, dst net.HardwareAddr, etherType uint16, msg []byte) {

	f := &ethernet.Frame{
		Destination: dst,
		Source:      src,
		EtherType:   ethernet.EtherType(etherType),
		Payload:     msg,
	}

	b, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal ethernet frame: %v", err)
	}

	addr := &raw.Addr{HardwareAddr: dst}
	if _, err := c.WriteTo(b, addr); err != nil {
		log.Fatalf("failed to send message: %v", err)
	}
}
