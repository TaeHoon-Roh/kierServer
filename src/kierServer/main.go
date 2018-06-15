package main

import (
	"fmt"
	"net"
	"os"
	"counter"
	"usepackbus"
	"io/ioutil"
)

type GlobalVal struct {
	Io_log []byte
	Io_last_tx bool
	Io_log_len int
}

func newGlobalVal() *GlobalVal {
	ngv := GlobalVal{}
	ngv.Io_last_tx = true
	ngv.Io_log = make([]byte, 16)
	ngv.Io_log_len = 0
	return &ngv
}

func GetLine(cnn *net.TCPConn) {
	read_buffer := make([]byte, 2048)
	read_index := 0
	getlinesig := uint16(0xAAAA)

	b := make([]byte, 1)
	unquote_next := false
	done := false

	counter := counter.Counter{}
	for done && counter.Elapsed() < 10000{
		result, err := cnn.Read(b)
		checkError(err)
		fmt.Println("result : " ,result)
		if b[0] == 0xBD{
			unquote_next = true
			continue
		}
		if b[0] != 0xBD{
			if unquote_next{
				b[0] -= 0x20
				unquote_next = false
			}
			getlinesig = CalcSigForByte(b[0],getlinesig)
			read_buffer[read_index] = b[0]
			read_index++
		} else if read_index >= 12 && getlinesig == 0{
			in_packet := usepackbus.Packet{}
			in_packet = in_packet.InitPacketParameter(read_buffer, read_index-2)
		}
	}
}
func CalcSigForByte(buff byte, seed uint16) uint16{
	rtn := seed
	j := rtn
	rtn = (rtn << 1) & 0x01FF
	if rtn >= 0x100 {
		rtn++
	}
	rtn = ((rtn + (j >> 8) + (uint16(buff))) & 0xFF) | (j << 8)
	return rtn

}
func main() {
	packet := usepackbus.Packet{}
	packet = packet.InitPacket()

	var GetCommandStatus int
	var NameList string
	var source_address uint16
	var dest_address uint16

	source_address = 4094
	dest_address = 4095

	cnn := tcptest()
	usepackbus.SetupDL(&packet, source_address, dest_address, &GetCommandStatus, &NameList, cnn)
	result, err := ioutil.ReadAll(cnn)
	checkError(err)

	fmt.Println(string(result))

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
