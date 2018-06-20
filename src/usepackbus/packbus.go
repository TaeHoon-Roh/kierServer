package usepackbus

import (
	"fmt"
	"unicode"
	"net"
	"counter"
)

type Packet struct {
	Storage     []byte
	Storage_len int
	Read_index  int

	Link_state       byte
	Dest_address     uint16
	Expect_more_code byte
	Priority         byte
	Src_address      uint16
	Hi_protocol_code uint16
	Message_type     byte
	Tran_no          byte

	Link_off_line byte
	Link_ring     byte
	Link_ready    byte
	Link_finished byte
	Link_pause    byte

	Expect_last    byte
	Expect_more    byte
	Expect_neutral byte
	Expect_reverse byte

	Pri_low        byte
	Pri_normal     byte
	Pri_high       byte
	Pri_extra_high byte

	Protocol_pakctrl byte
	Protocol_bmp5    byte
	Null_check       bool
}

func (p Packet) PrintPacket(pack *Packet) {
	fmt.Printf("\n Print Packet!!\n")
	fmt.Println("message type : ", pack.Message_type)
	fmt.Println("Storage Lenght : ", pack.Storage_len)
}

func (p Packet) InitPacket() (Packet) {
	fmt.Println("No Parameter Init Packet!!")
	p.Link_off_line = 8
	p.Link_ring = 9
	p.Link_ready = 10
	p.Link_finished = 11
	p.Link_pause = 12

	p.Expect_last = 0
	p.Expect_more = 1
	p.Expect_neutral = 2
	p.Expect_reverse = 3

	p.Pri_low = 0
	p.Pri_normal = 1
	p.Pri_high = 2
	p.Pri_extra_high = 3

	p.Protocol_pakctrl = 0
	p.Protocol_bmp5 = 1

	p.Link_state = p.Link_ready
	p.Expect_more_code = p.Expect_more
	p.Priority = p.Pri_high
	p.Src_address = 0
	p.Dest_address = 0
	p.Hi_protocol_code = uint16(p.Protocol_bmp5)
	p.Message_type = 0
	p.Tran_no = 0
	p.Storage_len = 0
	p.Read_index = p.Storage_len
	p.Null_check = false
	return p
}

func (p Packet) InitPacketParameter(buff []byte, len int) (Packet) {

	if len < 12 {
		//에러 메세지 넣어야함 throw new Exception("Invalid packet length");
		fmt.Println("Invalid packet Length")
	}

	word1 := int((int(buff[0]) << 8) | int(buff[1]))
	word2 := int((int(buff[2]) << 8) | int(buff[3]))
	word3 := int((int(buff[4]) << 8) | int(buff[5]))
	//word4 := int((int(buff[6]) << 8) | int(buff[7]))

	p.Link_state = byte((word1 & 0xF000) >> 12)
	p.Dest_address = uint16(word1 & 0x0FFF)
	p.Expect_more_code = byte((word2 & 0xC000) >> 14)
	p.Priority = byte((word2 & 0x0300) >> 12)
	p.Src_address = uint16(word2 & 0x0FFF)
	p.Hi_protocol_code = uint16((word3 & 0xF000) >> 12)

	p.Tran_no = 0
	p.Message_type = 0
	p.Storage_len = 0
	if len >= 9 {
		p.Message_type = buff[8]
		p.Tran_no = buff[9]
		p.Storage = make([]byte, len-10)
		for i := 10; i < len; i++ {
			p.Storage_len++
			p.Storage[i-10] = buff[i]
		}
	}
	p.Null_check = false
	return p
}

func (p Packet) Whats_left(packet *Packet) int {
	return (packet.Storage_len - packet.Read_index)
}

func (p Packet) reserve(pplen int, packet *Packet) {

	packet.Storage = make([]byte, 1)
	fmt.Println(len(packet.Storage))
	if packet.Storage == nil {
		packet.Storage = make([]byte, pplen)
	} else if len(packet.Storage) < pplen {
		temp := make([]byte, pplen*2)
		for i := 0; i < len(packet.Storage); i++ {
			temp[i] = packet.Storage[i]
		}
		packet.Storage = make([]byte, len(temp))
		copy(packet.Storage, temp)
	}

}

func (p Packet) add_bytes(buff []byte, buff_len int, packet *Packet) {
	p.reserve(packet.Storage_len+len(buff), packet)
	for i := 0; i < buff_len; i++ {
		packet.Storage[packet.Storage_len] = buff[i]
		packet.Storage_len++
	}

}

func (p Packet) add_byte(val byte, packet *Packet) {
	p.reserve(packet.Storage_len+1, packet)
	packet.Storage[packet.Storage_len] = val
	packet.Storage_len++
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
	rtn := make([]byte, packet.Storage_len+12)
	var i int

	fmt.Println("to_link_state_packet", p.Dest_address)
	rtn[0] = byte(((packet.Link_state << 4) | byte(((packet.Dest_address & 0xF00) >> 8))))
	rtn[1] = byte((packet.Dest_address & 0x00FF))
	rtn[2] = byte(((packet.Expect_more_code << 6) | byte((packet.Priority << 4)) | byte(((packet.Src_address & 0x0F00) >> 8))))
	rtn[3] = byte((packet.Src_address & 0x00FF));
	rtn[4] = byte((byte((packet.Hi_protocol_code << 4)) | byte(((packet.Dest_address & 0x0F00) >> 8))));
	rtn[5] = rtn[1];
	rtn[6] = byte(((packet.Src_address & 0x0F00) >> 8));
	rtn[7] = rtn[3];
	rtn[8] = packet.Message_type;
	rtn[9] = packet.Tran_no;
	for i = 0; i < packet.Storage_len; i++ {
		rtn[10+i] = packet.Storage[i]
	}

	var sig_null uint16 = calcSigNullifier(calcSigFor(rtn, int(packet.Storage_len+10), 0xAAAA))
	rtn[10+i] = byte(uint16(sig_null&0xFF00) >> 8)
	rtn[10+i+1] = byte(sig_null & 0x00FF)
	return rtn

}

func (p Packet) read_string(packet *Packet) string{
	var rtn string
	for (packet.Storage[packet.Read_index] != 0) && (packet.Read_index < packet.Storage_len){
		temp := rune(packet.Storage[packet.Read_index])
		rtn = rtn + string(temp)
		
	}

	return rtn
}

func (p Packet) read_byte(cnn *net.TCPConn, ioPackbus *IoPackbus) int {

	d := make([]byte, 1)
	cnn.Read(d)
	//fmt.Println("d : ", d[0])
	ioPackbus.Log_io(d[0], false, ioPackbus)
	return int(d[0])
}





type IoPackbus struct {
	Io_log     []byte
	Io_last_tx bool
	Io_log_len int
	flcount    int
	lcount     int
}

func (io IoPackbus) InitIoPackbus() IoPackbus {
	io.Io_log_len = 0
	io.Io_log = make([]byte, 16)
	io.Io_last_tx = true
	io.flcount = 0
	io.lcount = 0
	return io
}

func (io IoPackbus) Log_io(val byte, transmitted bool, packbus *IoPackbus) {
	fmt.Println("log_io count : ", packbus.lcount)
	packbus.lcount++
	if packbus.Io_last_tx != transmitted && packbus.Io_log_len > 0 {
		fmt.Println("check Log Io1")
		packbus.Flush_io_log("", packbus)

	}
	packbus.Io_log[packbus.Io_log_len] = val
	packbus.Io_log_len++
	packbus.Io_last_tx = transmitted
	if (packbus.Io_log_len == len(packbus.Io_log)) {
		fmt.Println("check Log Io2")
		packbus.Flush_io_log("", packbus)

	}
}

func (io IoPackbus) Flush_io_log(comment string, packbus *IoPackbus) {
	fmt.Println("flush log count : ", packbus.flcount)
	packbus.flcount++

	if packbus.Io_last_tx {
		fmt.Print("T ")
	} else {
		fmt.Print("R ")
	}

	for i := 0; i < packbus.Io_log_len; i += 1 {
		hex := io.Io_log[i] & 0x00FF
		/*if len(hex) == 1 {
			fmt.Print("0")
		}*/
		fmt.Printf("%.2x", hex)
		fmt.Print(" ")
	}
	for i := packbus.Io_log_len; i < len(packbus.Io_log); i += 1 {
		fmt.Printf("\t")
	}
	fmt.Print(" ")

	for i := 0; i < packbus.Io_log_len; i += 1 {
		ch := packbus.Io_log[i]
		if unicode.IsLetter(rune(ch)) {
			fmt.Printf("%c", ch)
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()
	if len(comment) > 0 {
		fmt.Println(comment)
	}
	packbus.Io_log_len = 0
}

func calcSigFor(buff []byte, len int, seed uint16) (uint16) {
	var j, n int
	rtn := seed
	fmt.Printf("check Len : %d\n", len)
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

func send_byte(val byte, cnn *net.TCPConn, packet *Packet, packbus *IoPackbus) {
	//fmt.Printf("Send Data!! : %x \n", val)
	//fmt.Println("check data ",packet.iopackbus.Io_log_len)
	packbus.Log_io(val, true, packbus)
	_, err := cnn.Write([]byte{val})
	if err != nil {
		fmt.Println(err)
	}
}

func SendPb(packet *Packet, cnn *net.TCPConn, packbus *IoPackbus) {

	var frame []byte = packet.to_link_state_packet(packet)
	//fmt.Printf("frame size : %d\n",len(frame) )
	send_byte(0xBD, cnn, packet, packbus)

	for i := 0; i < len(frame); i++ {
		if frame[i] == 0xBC || frame[i] == 0xBD {
			//fmt.Printf(" if : frame i : %d , %d \t", frame[i], i)
			send_byte(0xBC, cnn, packet, packbus)
			send_byte(byte(frame[i]+0x20), cnn, packet, packbus)
		} else {
			//fmt.Printf(" else : frame i : %d, %d \t", frame[i], i)
			send_byte(frame[i], cnn, packet, packbus)
		}
	}

	send_byte(0xBD, cnn, packet, packbus)
}

func SetupDL(packet *Packet, sa uint16, da uint16, gcs *int, nl *string, cnn *net.TCPConn, packbus *IoPackbus) {
	fmt.Println("Start setup DL")
	*gcs = 0
	*nl = "\000"
	GetCommand(packet, sa, da, nl)
	SendPb(packet, cnn, packbus)
	fmt.Println("End setup DL")
}

func GetCommand(packet *Packet, sa uint16, da uint16, nl *string) {
	CreateHeader(packet, sa, da)
	packet.Hi_protocol_code = uint16(packet.Protocol_pakctrl)
	packet.Message_type = 0x07
	packet.Tran_no = 0x07
	packet.add_string(nl, packet)
}

func CreateHeader(packet *Packet, sa uint16, da uint16) {
	packet.Src_address = sa
	packet.Dest_address = da
}

func GetLine(cnn *net.TCPConn, inpacket *Packet, ioPackbus *IoPackbus) {
	read_buffer := make([]byte, 2048)
	read_index := 0
	getlinesig := uint16(0xAAAA)

	var b byte
	unquote_next := false
	done := false

	counter := counter.Counter{}
	counter.InitCounter()
	for !done && counter.Elapsed() < 10000 {
		b = byte(inpacket.read_byte(cnn, ioPackbus))
		if b == 0xBC {
			fmt.Println("b == 0xBD")
			unquote_next = true
			continue
		}
		fmt.Printf("read_index : %d , getlinesig : %d , read_byte : %d\n",read_index,getlinesig,b)
		if b != 0xBD {
			if unquote_next {
				b -= 0x20
				unquote_next = false
			}
			getlinesig = CalcSigForByte(b, getlinesig)
			read_buffer[read_index] = b
			read_index++
		} else if read_index >= 12 && getlinesig == 0 {
			if read_index <= 12{
				ioPackbus.Flush_io_log("Invalid packet Length", ioPackbus)
				read_index = 0
				getlinesig = 0xAAAA
			} else {
				*inpacket = inpacket.InitPacketParameter(read_buffer, read_index-2)
				fmt.Println("check : ")
				str := "end of packet detected, size = " + string(inpacket.Whats_left(inpacket))
				ioPackbus.Flush_io_log(str, ioPackbus)
				done = true
			}
		} else {
			read_index = 0
			getlinesig = 0xAAAA
		}
	}
}



func CalcSigForByte(buff byte, seed uint16) uint16 {
	rtn := seed
	j := rtn
	rtn = (rtn << 1) & 0x01FF
	if rtn >= 0x100 {
		rtn++
	}
	rtn = ((rtn + (j >> 8) + (uint16(buff))) & 0xFF) | (j << 8)
	return rtn

}
