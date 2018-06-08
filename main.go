package main

import (
	"fmt"
	"net"
	"os"
)

type Packet struct {
	storage     []byte
	storage_len int
	read_index  int

	link_state       byte
	dest_address     int16
	expect_more_code byte
	priority         byte
	src_address      int16
	hi_protocol_code int16
	message_type     byte
	tran_no          byte

	link_off_line byte
	link_ring     byte
	link_ready    byte
	link_finished byte
	link_pause    byte

	expect_last    byte
	expect_more    byte
	expect_neutral byte
	expect_reverse byte

	pri_low        byte
	pri_normal     byte
	pri_high       byte
	pri_extra_high byte

	protocol_pakctrl byte
	protocol_bmp5    byte
}

func (p Packet) InitPacket() (Packet) {
	//p := Packet{}
	p.link_off_line = 8
	p.link_ring = 9
	p.link_ready = 10
	p.link_finished = 11
	p.link_pause = 12

	p.expect_last = 0
	p.expect_more = 1
	p.expect_neutral = 2
	p.expect_reverse = 3

	p.pri_low = 0
	p.pri_normal = 1
	p.pri_high = 2
	p.pri_extra_high = 3

	p.protocol_pakctrl = 0
	p.protocol_bmp5 = 1

	p.link_state = p.link_ready
	p.expect_more_code = p.expect_more
	p.priority = p.pri_high
	p.src_address = 0
	p.dest_address = 0
	p.hi_protocol_code = int16(p.protocol_bmp5)
	p.message_type = 0
	p.tran_no = 0
	p.storage_len = 0
	p.read_index = p.storage_len
	return p
}
func (p Packet) reserve(pplen int, packet *Packet) {

	packet.storage = make([]byte, 1)
	fmt.Println(len(packet.storage))
	if packet.storage == nil {
		packet.storage = make([]byte, pplen)
	} else if len(packet.storage) < pplen {
		temp := make([]byte, pplen*2)
		for i := 0; i < len(packet.storage); i++ {
			temp[i] = packet.storage[i]
		}
		packet.storage = make([]byte, len(temp))
		copy(packet.storage, temp)
	}

}
func (p Packet) add_bytes(buff []byte, buff_len int, packet *Packet) {
	fmt.Println("add_bytes")
	p.reserve(packet.storage_len+len(buff), packet)
	for i := 0; i < buff_len; i++ {
		packet.storage[packet.storage_len] = buff[i]
		packet.storage_len++
	}

}

func (p Packet) add_byte(val byte, packet *Packet) {
	fmt.Println("add_byte")
	p.reserve(packet.storage_len+1, packet)
	packet.storage[packet.storage_len] = val
	packet.storage_len++
}
func (p Packet) add_string(val *string, packet *Packet) {
	temp := []byte(*val)
	p.add_bytes(temp, len(temp), packet)
	if temp[len(temp)-1] == '0' && temp[len(temp)-2] == '\\' {
		fmt.Println("add_string : ", "check")
	} else {
		p.add_byte(0, packet)
	}
}

func requestHandler(c net.Conn) {
	data := make([]byte, 4096)

	for {
		n, err := c.Read(data)
		if err != nil {
			fmt.Print(err)
			return
		}

		fmt.Println(string(data[:n]))

	}
}

/*func makehttpHandle() {
	http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		fmt.Fprint(responseWriter, "Welcome to my website!")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static", http.StripPrefix("/static", fs))

	http.ListenAndServe(":80", nil)
}

*/


func calcSigFor(buff []byte, len int, seed uint16) (uint16) {
	var j, n int
	rtn := seed

	for n = 0 ; n < len ; n++{
		if n == 0 && buff[0] == 0xBD && len >1 {
			n = 1
			j = int(rtn)
			rtn = uint16((rtn << 1) & uint16(0x01FF))
			if int(rtn) >= int(0x100){
				rtn++
			}
			rtn = uint16(((rtn + uint16(j >> 8) + uint16(buff[0])) & uint16(0xFF)) | uint16(j << 8))
		}
	}

	return rtn
}
func calcSigNullifier(sig uint16) (uint16) {
	var new_seed uint16 = uint16 ((sig << 1) & uint16(0xFF))
	null1 := make([]byte, 1)
	var new_sig int = int(sig)

	if(new_seed >= 0x0100){
		new_seed++
	}
	null1[0] = byte(uint16(0x0100 - (new_seed + (sig >> 8))))
	new_sig = int(calcSigFor(null1, 1, sig))

	var null2 uint16

	new_seed = uint16(uint16(new_sig << 1) & uint16(0x01FF))
	if new_seed >= 0x0100 {
		new_seed++
	}
	null2 = uint16(uint16(0x0100 - (new_seed + uint16(new_sig >> 8))))

	rtn := uint16(null1[0])
	rtn <<= 8
	rtn += null2
	return rtn
}

func (p Packet) to_link_state_packet(packet *Packet) ([]byte) {
	rtn := make([]byte, packet.storage_len+12)
	var i int

	fmt.Println("to_link_state_packet",p.dest_address)
	rtn[0] = byte(((packet.link_state << 4) | byte(((packet.dest_address & 0xF00) >> 8))))
	rtn[1] = byte((packet.dest_address & 0x00FF))
	rtn[2] = byte(((packet.expect_more_code << 6) | byte((packet.priority<<4)) | byte(((packet.src_address&0x0F00)>>8))))
	rtn[3] = byte((packet.src_address & 0x00FF));
	rtn[4] = byte((byte((packet.hi_protocol_code << 4)) | byte(((packet.dest_address & 0x0F00) >> 8))));
	rtn[5] = rtn[1];
	rtn[6] = byte(((packet.src_address & 0x0F00) >> 8));
	rtn[7] = rtn[3];
	rtn[8] = packet.message_type;
	rtn[9] = packet.tran_no;

	for i = 0 ; i < packet.storage_len ; i++{
		rtn[10+i] = packet.storage[i]
	}
	
	var sig_null uint16 = calcSigNullifier(calcSigFor(rtn, packet.storage_len + 10, uint16(0xAAAA)))
	rtn[10 + i] = byte((sig_null & 0xFF00) >> 8)
	rtn[10 + i + 1] = byte(sig_null & 0x00FF)
	return rtn

}
func log_io(val byte, transmitted bool){

}
func send_byte(val byte, cnn *net.TCPConn)  {
	fmt.Printf("Send Data!! : %x \n" ,val)
	log_io(val, true)
	_, err := cnn.Write([]byte{val})
	if err != nil {
		fmt.Println(err)
	}
}
func SendPb(packet *Packet, cnn *net.TCPConn) {

	var frame []byte = packet.to_link_state_packet(packet)
	//fmt.Printf("frame size : %d\n",len(frame) )
	send_byte(0xBC, cnn)

	for i:= 0 ; i < len(frame) ; i++{
		if frame[i] == 0xBC || frame[i] == 0xBD {
			send_byte(0xBC, cnn)
			send_byte(byte(frame[i]+0x20), cnn)
		} else{
			send_byte(frame[i], cnn)
		}
	}

	send_byte(0xBD, cnn)
	//out_stream.flush()
}

func SetupDL(packet *Packet, sa int16, da int16, gcs *int, nl *string, cnn *net.TCPConn) {
	*gcs = 0
	*nl = "\\0"
	GetCommand(packet, sa, da, nl)
	SendPb(packet, cnn)
}

func GetCommand(packet *Packet, sa int16, da int16, nl *string) {
	CreateHeader(packet, sa, da)
	packet.hi_protocol_code = int16(packet.protocol_pakctrl)
	packet.message_type = 0x07
	packet.tran_no = 0x07
	packet.add_string(nl, packet)
}

func CreateHeader(packet *Packet, sa int16, da int16) {
	packet.src_address = sa
	packet.dest_address = da
}

func GetLine() {
/*	read_buffer := make([]byte, 2048)
	read_index := 0
	getlinesig := 0xAAAA

	var b byte
	unquote_net := false
	done := false
*/
}
func main() {
	packet := Packet{}
	packet = packet.InitPacket()

	var d uint16 = 'u'
	var dd = []byte(fmt.Sprintf("%x",d))
	fmt.Println(dd , " " , len(dd))

	var GetCommandStatus int
	var NameList string
	var source_address int16
	var dest_address int16

	source_address = 4094
	dest_address = 4095

	cnn := tcptest()

	SetupDL(&packet, source_address, dest_address, &GetCommandStatus, &NameList, cnn)


	/*result, err := ioutil.ReadAll(cnn)
	checkError(err)

	fmt.Println(string(result))*/

	os.Exit(0)
	/*shouldExit := false
	for {
		if shouldExit == true {
			break
		}

	}*/
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