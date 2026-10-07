package main

import (
	"io"
	"log"
	"net"
)

const (
	listenAddr  = ":8000"
	backendAddr = "localhost:9000"
)

func main() {

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	log.Printf("load balancer listening on %s", listenAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handleConnection(clientConn)

	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	backendConn, err := net.Dial("tcp", backendAddr)

	if err != nil {
		log.Printf("backend connection failed: %v", err)
		return
	}

	defer backendConn.Close()
	go func() {
		_, _ = io.Copy(backendConn, clientConn)
	}()
	_, _ = io.Copy(clientConn, backendConn)

}
