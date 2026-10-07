package main

import (
	"io"
	"log"
	"net"
	"sync"
)

const (
	listenAddr = ":8000"
)

var backends = []string{"localhost:9000", "localhost:9001", "localhost:9002"}
var counter int
var backendMu sync.Mutex

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
	var backendAddr string

	backendMu.Lock()
	defer backendMu.Unlock()
	counter = (counter + 1) % len(backends)
	backendAddr = backends[counter]
	backendMu.Unlock()

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
