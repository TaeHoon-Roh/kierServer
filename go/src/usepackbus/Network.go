package usepackbus

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func tcpSendBuffer(out_buffer *bufio.Writer, temp []byte) {

	fmt.Println("--Send Data!!--")

	out_buffer.WriteByte(0xBD)
	outnum, _ := out_buffer.Write(quote(temp))
	out_buffer.WriteByte(0xBD)
	fmt.Println("Socket Write size : ", outnum)
	out_buffer.Flush()
	for i := 0; i < len(temp); i++ {
		fmt.Printf("\\x%.2x", temp[i])
	}
	fmt.Println()
	//
	fmt.Println("--Send data end!!--")
	Send_count++

}

func tcpReadBuffer(in_buffer *bufio.Reader) ([]byte, int) {
	var pkt = make([]byte, 550)
	num, _ := in_buffer.Read(pkt)
	fmt.Println("Read data num is : ", num)
	result := uquote(pkt[1 : num-1])
	return result, len(result)
	/*pkt := make([]byte, 0)
	for {
		b, err := in_buffer.ReadByte()
		if err != nil {
			break
		} else {
			pkt = append(pkt,b)
		}

	}

	if len(pkt) == 0{
		fmt.Println("Don't Read")
	} else {
		result := uquote(pkt[1:len(pkt)-1])
		return result, len(result)
	}
	return nil, 0*/
}

func ConnectDevice(IpAddress string) (*net.TCPConn) {
	service := IpAddress

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	checkError(err)

	return conn
}