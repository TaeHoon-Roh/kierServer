package usepackbus

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func tcpSendBuffer(out_buffer *bufio.Writer, temp []byte) {

	fmt.Println("--Send Data!!--")
	buff := quote(temp)
	time.Sleep(time.Microsecond)
	fmt.Println("temp size : ", len(buff))
	fmt.Println("temp check : ", temp)
	out_buffer.WriteByte(0xBD)
	outnum, _ := out_buffer.Write(buff)
	out_buffer.WriteByte(0xBD)
	fmt.Println("Socket Write size : ", outnum)
	out_buffer.Flush()
	out_buffer.Flush()
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
	if num == 0{
		return nil, -999
	}
	result := uquote(pkt[1 : num-1])
	for i := 0 ; i < len(result) ; i++ {
		fmt.Printf("\\%0.2x", result[i])
	}
	fmt.Println()
	return result, len(result)
}

func ConnectDevice(IpAddress string) (*net.TCPConn) {
	service := IpAddress

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	checkError(err)

	return conn
}