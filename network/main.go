package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var localPort = flag.Int("p", 25565, "Local Port to listen on")
var remoteHost = flag.String("r", "localhost:25566", "Remote Server address host:port")

func main() {
	flag.Parse()
	w := wrapper{
		0,
		*localPort,
		*remoteHost,
	}
	w.start()
}

type wrapper struct {
	ids        int
	localPort  int
	remoteHost string
}

func (w *wrapper) start() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", w.localPort))
	if err != nil {
		fmt.Printf("failed to listen: %s\n", err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("[*] Listening...")

	for {
		clientConn, err := l.Accept()
		if err != nil {
			fmt.Printf("failed to accept: %s\n", err)
			break
		}
		fmt.Printf("[*] Accepted from: %s\n", clientConn.RemoteAddr())
		go w.handleConnection(clientConn)
	}
}

func (w *wrapper) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()
	serverConn, err := net.Dial("tcp", w.remoteHost)
	if err != nil {
		fmt.Printf("failed to dial: %s\n", err)
		return
	}
	defer serverConn.Close()
	fmt.Printf("[*][%d] Connected to server: %s\n", w.ids, serverConn.RemoteAddr())
	id := w.ids
	w.ids++
	go w.handleServerMessages(serverConn, clientConn, id)
	w.handleClientMessages(serverConn, clientConn, id)
}

func (w *wrapper) handleServerMessages(serverConn, clientConn net.Conn, id int) {
	for {
		data := make([]byte, 2048)
		n, err := serverConn.Read(data)
		if err != nil && err != io.EOF {
			fmt.Printf("failed to read from server: %s\n", err)
			break
		}
		if n > 0 {
			fmt.Printf("From Server [%d]:\n%s\n", id, hex.Dump(data[:n]))
			clientConn.Write(data[:n])
		}
	}
}

func (w *wrapper) handleClientMessages(serverConn, clientConn net.Conn, id int) {
	for {
		data := make([]byte, 2048)
		n, err := clientConn.Read(data)
		if err != nil && err == io.EOF {
			fmt.Printf("failed to read from client: %s\n", err)
			break
		}
		if n > 0 {
			fmt.Printf("From Client [%d]:\n%s\n", id, hex.Dump(data[:n]))
			serverConn.Write(data[:n])
		}
	}
}
