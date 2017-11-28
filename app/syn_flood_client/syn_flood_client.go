package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"net"
	"syscall"

	"github.com/iyidan/tcpdemo/common"
)

func main() {
	for i := 0; i < 1000; i++ {
		sockfd := sendSyn(net.ParseIP("192.168.36.135").To4(), net.ParseIP("192.168.36.1").To4(), uint16(i+30000), 18001)
		syscall.Shutdown(sockfd, syscall.SHUT_RDWR)
	}
}

func sendSyn(srcIP, dstIP net.IP, srcPort, dstPort uint16) (sockfd int) {
	// tcp伪首部
	tcpPsdHeader := common.TCPPsdHeader{
		SrcIP:    uint32(srcIP[0])<<24 + uint32(srcIP[1])<<16 + uint32(srcIP[2])<<8 + uint32(srcIP[3]),
		DstIP:    uint32(dstIP[0])<<24 + uint32(dstIP[1])<<16 + uint32(dstIP[2])<<8 + uint32(dstIP[3]),
		Reversed: 0x00,
		Protocol: syscall.IPPROTO_TCP,
		TCPLen:   20, // syn固定20字节
	}
	// tcp首部
	tcpHeader := common.TCPHeader{
		SrcPort:   srcPort,
		DstPort:   dstPort,
		SeqNum:    1,
		AckNum:    0,
		Offset:    5 << 4, // 前4字节
		Flag:      1 << 1, // 后6字节
		Window:    1500,
		Checksum:  0,
		UrgentPtr: 0,
	}

	var (
		buf bytes.Buffer
	)
	binary.Write(&buf, binary.BigEndian, tcpPsdHeader)
	binary.Write(&buf, binary.BigEndian, tcpHeader)
	psdPackage := buf.Bytes()
	tcpHeader.Checksum = common.Checksum(psdPackage)
	log.Printf("srcPort: %d, psd Package:\n%v\n", srcPort, hex.Dump(psdPackage))

	buf.Reset()
	binary.Write(&buf, binary.BigEndian, tcpHeader)

	synPackage := buf.Bytes()
	log.Printf("srcPort: %d, syn package:\n%v\n", srcPort, hex.Dump(synPackage))

	// 原始套接字
	var (
		addr syscall.SockaddrInet4
		err  error
	)
	if sockfd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP); err != nil {
		log.Println("Socket() error: ", err.Error())
		return
	}
	addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3] = 192, 168, 36, 1
	addr.Port = int(tcpHeader.DstPort)
	if err = syscall.Sendto(sockfd, synPackage, 0, &addr); err != nil {
		log.Println("Sendto() error: ", err.Error())
		return
	}
	log.Println("Send success!")
	return
}
