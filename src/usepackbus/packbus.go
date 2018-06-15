package usepackbus

import (
	"fmt"
	"net"
)

type Packet struct {
	iopackbus IoPackbus
	storage     []byte
	storage_len int
	read_index  int

	link_state       byte
	dest_address     uint16
	expect_more_code byte
	priority         byte
	src_address      uint16
	hi_protocol_code uint16
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
	fmt.Println("No Parameter Init Packet!!")
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
	p.hi_protocol_code = uint16(p.protocol_bmp5)
	p.message_type = 0
	p.tran_no = 0
	p.storage_len = 0
	p.read_index = p.storage_len
	p.iopackbus = IoPackbus{}
	p.iopackbus.InitIoPackbus()
	return p
}

func (p Packet) InitPacketParameter(buff []byte, len int) (Packet) {

	if len < 12{
		//에러 메세지 넣어야함 throw new Exception("Invalid packet length");
	}

	word1 := int((int(buff[0]) << 8) | int(buff[1]))
	word2 := int((int(buff[2]) << 8) | int(buff[3]))
	word3 := int((int(buff[4]) << 8) | int(buff[5]))
	//word4 := int((int(buff[6]) << 8) | int(buff[7]))

	p.link_state = byte((word1 & 0xF000) >> 12)
	p.dest_address = uint16(word1 & 0x0FFF)
	p.expect_more_code = byte((word2 & 0xC000) >> 14)
	p.priority = byte((word2 & 0x0300) >> 12)
	p.src_address = uint16(word2 & 0x0FFF)
	p.hi_protocol_code = uint16((word3 & 0xF000) >> 12)

	p.tran_no = 0
	p.message_type = 0
	p. storage_len = 0
	if len >= 9{
		p.message_type = buff[8]
		p.tran_no = buff[9]
		p.storage = make([]byte, len-10)
		for i := 10 ; i < len ; i++{
			p.storage_len++
			p.storage[i-10] = buff[i]
		}
	}
	return p
}

func (p Packet) Whats_left(packet Packet) int  {
	return (packet.storage_len - packet.read_index)
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
	fmt.Printf("temp : %d temp length : %d\n", temp, len(temp))
	p.add_bytes(temp, len(temp), packet)
	if temp[len(temp)-1] != '\000' {
		fmt.Println("add_string : ", "check")
		p.add_byte(0, packet)
	}
}

func (p Packet) to_link_state_packet(packet *Packet) ([]byte) {
	rtn := make([]byte, packet.storage_len+12)
	var i int

	fmt.Println("to_link_state_packet", p.dest_address)
	rtn[0] = byte(((packet.link_state << 4) | byte(((packet.dest_address & 0xF00) >> 8))))
	rtn[1] = byte((packet.dest_address & 0x00FF))
	rtn[2] = byte(((packet.expect_more_code << 6) | byte((packet.priority << 4)) | byte(((packet.src_address & 0x0F00) >> 8))))
	rtn[3] = byte((packet.src_address & 0x00FF));
	rtn[4] = byte((byte((packet.hi_protocol_code << 4)) | byte(((packet.dest_address & 0x0F00) >> 8))));
	rtn[5] = rtn[1];
	rtn[6] = byte(((packet.src_address & 0x0F00) >> 8));
	rtn[7] = rtn[3];
	rtn[8] = packet.message_type;
	rtn[9] = packet.tran_no;
	fmt.Printf("packet.storage_len : %d\n", packet.storage_len)
	for i = 0; i < packet.storage_len; i++ {
		rtn[10+i] = packet.storage[i]
	}

	var sig_null uint16 = calcSigNullifier(calcSigFor(rtn, int(packet.storage_len+10), 0xAAAA))
	rtn[10+i] = byte(uint16(sig_null & 0xFF00) >> 8)
	rtn[10+i+1] = byte(sig_null & 0x00FF)
	return rtn

}

func calcSigFor(buff []byte, len int, seed uint16) (uint16) {
	var j, n int
	rtn := seed
	fmt.Printf("check Len : %d\n",len)
	for n = 0; n < len; n++ {
		if n == 0 && buff[0] == 0xBD && len > 1 {
			n = 1
		}
		j = int(rtn)
		rtn = uint16((rtn << 1) & 0x01FF)
		if int(rtn) >= int(0x100) {
			rtn++
		}
		rtn = uint16(((rtn + uint16(j>>8) + uint16(buff[n])) & uint16(0xFF)) | uint16(j<<8))
	}
	return rtn
}
func calcSigNullifier(sig uint16) (uint16) {
	fmt.Printf("Sig : %d\n", sig)
	var new_seed uint16 = uint16((sig << 1) & uint16(0x1FF))
	null1 := make([]byte, 1)
	var new_sig int = int(sig)

	if (new_seed >= 0x0100) {
		new_seed++
	}
	null1[0] = byte(uint16(0x0100 - (new_seed + (sig >> 8))))
	new_sig = int(calcSigFor(null1, 1, uint16(sig)))

	var null2 uint16

	new_seed = uint16(uint16(new_sig<<1) & uint16(0x01FF))
	if new_seed >= 0x0100 {
		new_seed++
	}
	null2 = uint16(uint16(0x0100 - (new_seed + uint16(new_sig>>8))))

	rtn := uint16(null1[0])
	rtn <<= 8
	rtn += null2
	return rtn
}



func send_byte(val byte, cnn *net.TCPConn, packet *Packet) {
	//fmt.Printf("Send Data!! : %x \n", val)
	//fmt.Println("check data ",packet.iopackbus.Io_log_len)
	packet.iopackbus.Log_io(val, true)
	_, err := cnn.Write([]byte{val})
	if err != nil {
		fmt.Println(err)
	}
}

func SendPb(packet *Packet, cnn *net.TCPConn) {

	var frame []byte = packet.to_link_state_packet(packet)
	//fmt.Printf("frame size : %d\n",len(frame) )
	send_byte(0xBD, cnn, packet)

	for i := 0; i < len(frame); i++ {
		if frame[i] == 0xBC || frame[i] == 0xBD {
			//fmt.Printf(" if : frame i : %d , %d \t", frame[i], i)
			send_byte(0xBC, cnn, packet)
			send_byte(byte(frame[i]+0x20), cnn, packet)
		} else {
			//fmt.Printf(" else : frame i : %d, %d \t", frame[i], i)
			send_byte(frame[i], cnn, packet)
		}
	}

	send_byte(0xBD, cnn, packet)
}

func SetupDL(packet *Packet, sa uint16, da uint16, gcs *int, nl *string, cnn *net.TCPConn) {
	*gcs = 0
	*nl = "\000"
	GetCommand(packet, sa, da, nl)
	SendPb(packet, cnn)
}

func GetCommand(packet *Packet, sa uint16, da uint16, nl *string) {
	CreateHeader(packet, sa, da)
	packet.hi_protocol_code = uint16(packet.protocol_pakctrl)
	packet.message_type = 0x07
	packet.tran_no = 0x07
	packet.add_string(nl, packet)
}

func CreateHeader(packet *Packet, sa uint16, da uint16) {
	packet.src_address = sa
	packet.dest_address = da
}

type IoPackbus struct {
	Io_log []byte
	Io_last_tx bool
	Io_log_len int
}

func (io *IoPackbus) InitIoPackbus() {
	//nip := IoPackbus{}
	io.Io_log_len = 0
	io.Io_log = make([]byte, 16)
	io.Io_last_tx = true
}


func (io *IoPackbus) Log_io(val byte, transmitted bool) {
	if io.Io_last_tx != transmitted && io.Io_log_len > 0{
		fmt.Println("check Log Io1")
	}
	io.Io_log[io.Io_log_len] = val
	io.Io_log_len++
	io.Io_last_tx = transmitted
	if(io.Io_log_len == len(io.Io_log)){
		fmt.Println("check Log Io2")
	}
}

func (io *IoPackbus) Flush_io_log(comment string)  {

	if io.Io_last_tx{
		fmt.Print("T ")
	} else {
		fmt.Print("R ")
	}

	for i := 0 ; i < io.Io_log_len ; i += 1{
		hex := string(io.Io_log[i] & 0x00FF)
		if len(hex) == 1 {
			fmt.Print("0")
		}
		fmt.Print(hex)
		fmt.Print(" ")
	}

	for i:= io.Io_log_len ; i < len(io.Io_log) ; i += 1{
		fmt.Printf("\t")
	}
	fmt.Print(" ")

	for i:= 0 ; i < io.Io_log_len ; i += 1{
		ch := uint16(io.Io_log[i])
		
	}


}