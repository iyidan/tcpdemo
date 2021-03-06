package main

import (
	"log"
	"net"
	"time"
)

func main() {
	//conn, err := net.Dial("tcp", "192.168.36.1:18001")
	conn, err := net.Dial("tcp", "192.168.36.135:18001")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	// tcpConn := conn.(*net.TCPConn)
	// err = tcpConn.SetKeepAlive(true)
	// if err != nil {
	// 	log.Fatal("tcpConn.SetKeepAlive:", err)
	// }
	// err = tcpConn.SetKeepAlivePeriod(time.Second * 1)
	// if err != nil {
	// 	log.Fatal("tcpConn.SetKeepAlivePeriod:", err)
	// }

	if err := conn.(*net.TCPConn).SetReadBuffer(100); err != nil {
		log.Fatal("TCPConn.SetReadBuffer:", err)
	}

	for {
		data := make([]byte, 1)
		n, err := conn.Read(data)
		log.Printf("conn.Read, n: %d, data: %s, err: %v", n, data, err)
		time.Sleep(time.Second * 5)
	}
}
