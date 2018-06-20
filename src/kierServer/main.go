package main

import (
	"fmt"
	"net"
	"os"
	"usepackbus"
)

var ioPackbus usepackbus.IoPackbus

func main() {
	//main_tcp_packbus()

}

func main_tcp_packbus(){

	ioPackbus := usepackbus.IoPackbus{}
	ioPackbus = ioPackbus.InitIoPackbus()
	packet := usepackbus.Packet{}
	packet = packet.InitPacket()
	fmt.Println("IOpackbus : " ,ioPackbus.Io_log_len)


	var GetCommandStatus int
	var NameList string
	var source_address uint16
	var dest_address uint16

	source_address = 4094
	dest_address = 4095

	cnn := tcptest()
	usepackbus.SetupDL(&packet, source_address, dest_address, &GetCommandStatus, &NameList, cnn, &ioPackbus)


	inpacket := usepackbus.Packet{}
	usepackbus.GetLine(cnn, &inpacket, &ioPackbus)

	switch (int(inpacket.Message_type)&0xFF) {
	case 0x97:
		fmt.Println("Received clock check response")
		break
	case 0x89:
		fmt.Println("Response for Collect Data Transaction")
		break
	case 0x87:
		fmt.Println("Response for Command format")
		break
	case 0x9a:
		break
	case 0x9d:
		fmt.Println("Response for Upload Command")
		break
	}

	cnn.Close()

	os.Exit(0)

}

func tcptest() (*net.TCPConn) {
	service := "147.46.138.187:6785"

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	return conn
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s", err.Error())
		os.Exit(1)
	}

}