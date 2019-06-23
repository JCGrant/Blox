package proxy

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pkg/errors"
)

var localPort = 25565
var remoteHost = "localhost:25566"

// Proxy sits in between the Minecraft server and the Minecraft client
type Proxy struct {
	ids        int
	localPort  int
	remoteHost string
}

// New creates a Proxy
func New() *Proxy {
	return &Proxy{
		0,
		localPort,
		remoteHost,
	}
}

// Start starts the Proxy
func (p *Proxy) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", p.localPort))
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}
	defer l.Close()

	log.Println("[*] Listening...")

	for {
		clientConn, err := l.Accept()
		if err != nil {
			log.Printf("failed to accept: %s\n", err)
			break
		}
		log.Printf("[*] Accepted from: %s\n", clientConn.RemoteAddr())
		go p.handleConnection(clientConn)
	}
	return nil
}

func (p *Proxy) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()
	serverConn, err := net.Dial("tcp", p.remoteHost)
	if err != nil {
		log.Printf("failed to dial: %s\n", err)
		return
	}
	defer serverConn.Close()
	log.Printf("[*][%d] Connected to server: %s\n", p.ids, serverConn.RemoteAddr())
	id := p.ids
	p.ids++
	go p.handleServerMessages(serverConn, clientConn, id)
	p.handleClientMessages(serverConn, clientConn, id)
}

func (p *Proxy) handleServerMessages(serverConn, clientConn net.Conn, id int) {
	for {
		data := make([]byte, 2048)
		n, err := serverConn.Read(data)
		if err != nil && err != io.EOF {
			log.Printf("failed to read from server: %s\n", err)
			break
		}
		if n > 0 {
			log.Printf("From Server [%d]:\n%s\n", id, hex.Dump(data[:n]))
			clientConn.Write(data[:n])
		}
	}
}

func (p *Proxy) handleClientMessages(serverConn, clientConn net.Conn, id int) {
	for {
		data := make([]byte, 2048)
		n, err := clientConn.Read(data)
		if err != nil && err == io.EOF {
			log.Printf("failed to read from client: %s\n", err)
			break
		}
		if n > 0 {
			log.Printf("From Client [%d]:\n%s\n", id, hex.Dump(data[:n]))
			serverConn.Write(data[:n])
		}
	}
}
