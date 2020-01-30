package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
)

type Endpoint struct {
	User string
	Host string
	Port int
}

func NewEndpoint(user, host string, port int) Endpoint {
	return Endpoint{User: user, Host: host, Port: port}
}

func (endpoint Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type Tunnel struct {
	LocalEndpoint  Endpoint
	bastionClient  *ssh.Client
	remoteEndpoint Endpoint
	localListener  net.Listener
}

func NewTunnel(bastionEndpoint, remoteEndpoint Endpoint, auth ssh.AuthMethod) (Tunnel, error) {

	localListener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return Tunnel{}, fmt.Errorf("ssh tunnel: local listen: %v", err)
	}
	localEndpoint := Endpoint{Host: "localhost", Port: localListener.Addr().(*net.TCPAddr).Port}
	log.Printf("ssh tunnel: local listen on %s", localEndpoint)

	bastionClient, err := newBastionClient(bastionEndpoint, auth)
	if err != nil {
		return Tunnel{}, fmt.Errorf("ssh tunnel: bastion dial: %v", err)
	}
	log.Printf("ssh tunnel: connected to bastion %s", bastionEndpoint.String())

	return Tunnel{
		LocalEndpoint:  localEndpoint,
		localListener:  localListener,
		bastionClient:  bastionClient,
		remoteEndpoint: remoteEndpoint,
	}, nil
}

func newBastionClient(bastionEndpoint Endpoint, auth ssh.AuthMethod) (*ssh.Client, error) {

	sshConfig := &ssh.ClientConfig{
		User: bastionEndpoint.User,
		Auth: []ssh.AuthMethod{auth},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Always accept key.
			return nil
		},
	}
	return ssh.Dial("tcp", bastionEndpoint.String(), sshConfig)
}

// blocking call, start in a go routine
func (tunnel Tunnel) Start() {

	if tunnel.localListener == nil {
		log.Printf("tunnel local listener is nil")
		return
	}

	for {
		localConn, err := tunnel.localListener.Accept()
		if err != nil {
			log.Printf("ssh tunnel: local listener accept: %v", err)
			return
		}
		log.Printf("ssh tunnel: accepted connection")

		remoteConn, err := tunnel.remoteDial()
		if err != nil {
			log.Printf("ssh tunnel: local listener accept: %v", err)
			return
		}
		forward(localConn, remoteConn)
	}
}

func (tunnel Tunnel) remoteDial() (net.Conn, error) {

	remoteConn, err := tunnel.bastionClient.Dial("tcp", tunnel.remoteEndpoint.String())
	if err != nil {
		return nil, fmt.Errorf("ssh tunnel: remote dial error: %v", err)
	}
	log.Printf("ssh tunnel: connected to remote %s", tunnel.remoteEndpoint.String())
	return remoteConn, nil
}

func forward(localConn, remoteConn net.Conn) {

	log.Printf("ssh tunnel: forwarding local <-> remote")
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			log.Printf("ssh tunnel: io.Copy error: %v", err)
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}
