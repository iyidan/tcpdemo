package common

// Checksum 计算首部校验和
// 1. 把校验和字段置为0；
// 2. 对IP头部中的每16bit进行二进制求和；
// 3. 如果和的高16bit不为0，则将和的高16bit和低16bit反复相加，
// 直到和的高16bit为0，从而获得一个16bit的值；
// 4. 将该16bit的值取反，存入校验和字段。
func Checksum(msg []byte) uint16 {
	sum := uint32(0)
	for n := 0; n < len(msg); n += 2 {
		sum += uint32(msg[n])<<8 + uint32(msg[n+1])
	}
	for sum>>16 != 0 {
		sum = sum>>16 + sum&0xffff
	}
	return uint16(^sum)
}

// TCPHeader TCP报文头部格式，不包括可变长头部
type TCPHeader struct {
	// 源端口， 16bit
	SrcPort uint16
	// 目标端口，16bit
	DstPort uint16
	// 数据序号，32bits，TCP 连接中传送的数据流中的每一个字节都编上一个序号。
	// 序号字段的值则指的是本报文段所发送的数据的第一个字节的序号
	SeqNum uint32
	// 确认号，32bits，期望收到对方的下一个报文段的数据的第一个字节的序号。
	AckNum uint32
	// 数据偏移，4bits，单位为4字节，它指出报文数据距TCP报头的起始处有多远(TCP报文头长度)。
	// 保留字段 6bits，保留今后使用，目前置0
	Offset uint8
	// URG：紧急比特，1bit，当 URG=1 时，表明紧急指针字段有效。它告诉系统此报文段中有紧急数据，应尽快传送(相当于高优先级的数据)
	// ACK：确认比特，1bit，只有当 ACK=1时确认号字段才有效。当 ACK=0 时，确认号无效
	// PSH：推送比特，1bit，接收方 TCP 收到推送比特置1的报文段，就尽快地交付给接收应用进程，而不再等到整个缓存都填满了后再向上交付
	// RST：复位比特，1bit，当RST=1时，表明TCP连接中出现严重差错(如由于主机崩溃或其他原因)，必须释放连接，然后再重新建立运输连接
	// SYN：同步比特，1bit，同步比特 SYN 置为 1，就表示这是一个连接请求或连接接受报文
	// FIN：终止比特，1bit，用来释放一个连接。当FIN=1 时，表明此报文段的发送端的数据已发送完毕，并要求释放运输连接
	// Flag前2bit为保留字段，目前置0
	Flag uint8
	// 窗口大小，16bits，窗口字段用来控制对方发送的数据量，单位为字节。
	// TCP 连接的一端根据设置的缓存空间大小确定自己的接收窗口大小，
	// 然后通知对方以确定对方的发送窗口的上限
	Window uint16
	// 检验和，16bits，检验和字段检验的范围包括首部和数据这两部分。
	// 在计算检验和时，要在 TCP 报文段的前面加上 12 字节的伪首部
	Checksum uint16
	// 紧急指针字段，16bits，紧急指针指出在本报文段中的紧急数据的最后一个字节的序号
	UrgentPtr uint16
}

// TCPPsdHeader 伪首部共有12字节，包含如下信息：
// 源IP地址(32bit)
// 目的IP地址(32bit)
// 保留字节，置0(8bit)
// 传输层协议号，TCP是6(8bit)
// TCP报文长度(报头+数据)(16bit)
type TCPPsdHeader struct {
	SrcIP    uint32
	DstIP    uint32
	Reversed uint8
	Protocol uint8
	TCPLen   uint16
}
