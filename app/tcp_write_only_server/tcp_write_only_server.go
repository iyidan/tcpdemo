package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"syscall"
)

func main() {
	addr := ":18001"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("net.Listen:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("listener.Accept:", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	connID := fmt.Sprintf("%s <-> %s:", conn.RemoteAddr().String(), conn.LocalAddr().String())
	log.Println(connID, "connected")

	// n, err := handleConnDefault(conn)
	n, err := handleConnTCPCORK(conn)

	log.Printf("%s conn write err, n: %d, err: %v", connID, n, err)
}

func handleConnDefault(conn net.Conn) (int, error) {
	for {
		data := bytes.Repeat([]byte("a"), 1024*1024)
		if n, err := conn.Write(data); err != nil {
			return n, err
		}
		//time.Sleep(time.Second * 1)
	}
}

func handleConnTCPCORK(conn net.Conn) (n int, err error) {
	err = conn.(*net.TCPConn).SetNoDelay(false)
	if err != nil {
		return
	}
	rawConn, err := conn.(*net.TCPConn).SyscallConn()
	if err != nil {
		return
	}
	err = rawConn.Control(func(fd uintptr) {
		val, err := syscall.GetsockoptInt(int(fd), 0x6, 0x3)
		log.Printf("GetsockoptInt, val: %d, err: %v", val, err)
		syscall.SetsockoptInt(int(fd), 0x6, 0x3, 1)
	})
	if err != nil {
		return
	}

	for {
		data := bytes.Repeat([]byte("a"), 1024*1024)
		if n, err = conn.Write(data); err != nil {
			return
		}
		//time.Sleep(time.Second * 1)
	}
}
