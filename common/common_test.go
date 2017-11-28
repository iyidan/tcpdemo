package common

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"testing"
)

func TestIPChecksum(t *testing.T) {
	headerCase := []uint16{
		0x4500, 0x0047,
		0x1fd6, 0x0000,
		0x0111, 0xf7be,
		0xc0a8, 0x006d,
		0xe000, 0x00fc,
	}

	msg := make([]byte, 20)
	for i, v := range headerCase {
		if i == 5 {
			binary.BigEndian.PutUint16(msg[i*2:], uint16(0x0000))
			continue
		}
		binary.BigEndian.PutUint16(msg[i*2:], v)
	}
	t.Logf("msg:\n%v", hex.Dump(msg))
	sum := Checksum(msg)
	t.Logf("sum: %x", sum)
	if sum != headerCase[5] {
		t.Fatal()
	}
}

func TestTCPChecksum(t *testing.T) {
	srcIP := net.IPv4(0xc0, 0xa8, 0x24, 0x87).To4()
	dstIP := net.IPv4(0xc0, 0xa8, 0x24, 0x01).To4()
	t.Log("srcIP:", srcIP.String())
	t.Log("dstIP:", dstIP.String())
	t.Logf("srcIP-Hex: %x", uint32(srcIP[0])<<24+uint32(srcIP[1])<<16+uint32(srcIP[2])<<8+uint32(srcIP[3]))
	t.Logf("srcIP-Hex: %x", uint32(dstIP[0])<<24+uint32(dstIP[1])<<16+uint32(dstIP[2])<<8+uint32(dstIP[3]))

	headerCase := []uint16{
		0xc0a8, 0x2487, // srcIP: 32bit
		0xc0a8, 0x2401, // dstIP: 32bit
		0x0006, 0x0014, // reversed, protocol, tcp package length: 8bit+8bit+16bit
		0x51a3, 0x4651,
		0x0000, 0x0001,
		0x0000, 0x0000,
		0x5002, 0x05dc,
		0x4838, 0x0000,
	}

	msg := make([]byte, len(headerCase)*2)
	for i, v := range headerCase {
		if i == len(headerCase)-2 {
			binary.BigEndian.PutUint16(msg[i*2:], uint16(0x0000))
			continue
		}
		binary.BigEndian.PutUint16(msg[i*2:], v)
	}
	t.Logf("msg:\n%v", hex.Dump(msg))
	sum := Checksum(msg)
	t.Logf("sum: %x", sum)
	if sum != headerCase[len(headerCase)-2] {
		t.Fatal()
	}
}
