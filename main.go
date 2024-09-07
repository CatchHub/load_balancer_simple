package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	counter int
	// TODO configurable
	listenAdr = "localhost:8080"
	// TODO configurable
	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)

func main() {
	listener, err := net.Listen("tcp", listenAdr)

	if err != nil {
		log.Fatal("failed to listen: %s", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err)
		}

		backend := chooseBackend()
		fmt.Println("counter=%d, backend=%s\n", counter, backend)
		go func() {
			err := proxy(backend, conn)
			if err != nil {
				log.Printf("WARNING: proxying failed: %v", err)
			}
		}()
	}
}

func proxy(backend string, c net.Conn) error {
	bc, err := net.Dial("tcp", backend)
	if err != nil {
		return fmt.Errorf("failed connect to backend %s: %v", backend, err)
	}

	// c => bc
	go io.Copy(bc, c)
	// bc => c
	go io.Copy(c, bc)
	return nil
}

func chooseBackend() string {
	s := server[counter%len(server)]
	counter++
	return s
}
