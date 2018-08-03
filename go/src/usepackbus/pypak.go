package usepackbus

import (
	"fmt"
	"encoding/binary"
	"strings"
	"encoding/hex"
)

type DataTypeDef struct {
	Name   string
	Code   byte
	Format string
	Size   byte
}


func InitDataTypeDef() ([]DataTypeDef){

	datatype := make([]DataTypeDef, 27)

	datatype[0] = DataTypeDef{"Byte", 1, "B", 1}
	datatype[1] = DataTypeDef{"UInt2", 2, ">H", 2}
	datatype[2] = DataTypeDef{"UInt4", 3, ">L", 4}
	datatype[3] = DataTypeDef{"Int1", 4, "b", 1}
	datatype[4] = DataTypeDef{"Int2", 5, ">h", 2}
	datatype[5] = DataTypeDef{"Int4", 6, ">l", 4}
	datatype[6] = DataTypeDef{"FP2", 7, ">H", 2}
	datatype[7] = DataTypeDef{"FP3", 15, "3c", 3}
	datatype[8] = DataTypeDef{"FP4", 8, "4c", 4}
	datatype[9] = DataTypeDef{"IEEE4B", 9, ">f", 4}
	datatype[10] = DataTypeDef{"IEEE8B", 18, ">d", 8}
	datatype[11] = DataTypeDef{"Bool8", 17, "B", 1}
	datatype[12] = DataTypeDef{"Bool", 10, "B", 1}
	datatype[13] = DataTypeDef{"Bool2", 27, ">H", 2}
	datatype[14] = DataTypeDef{"Bool4", 28, ">L", 4}
	datatype[15] = DataTypeDef{"Sec", 12, ">l", 4}
	datatype[16] = DataTypeDef{"USec", 13, "6c", 6}
	datatype[17] = DataTypeDef{"NSec", 14, ">2l", 8}
	datatype[18] = DataTypeDef{"ASCII", 11, "s", 0}
	datatype[19] = DataTypeDef{"ASCIIZ", 16, "s", 0}
	datatype[20] = DataTypeDef{"Short", 19, "<h", 2}
	datatype[21] = DataTypeDef{"Long", 20, "<l", 4}
	datatype[22] = DataTypeDef{"UShort", 21, "<H", 2}
	datatype[23] = DataTypeDef{"ULong", 22, "<L", 4}
	datatype[24] = DataTypeDef{"IEEE4L", 24, "<f", 4}
	datatype[25] = DataTypeDef{"IEEE8L", 25, "<d", 8}
	datatype[26] = DataTypeDef{"SecNano", 23, "<2l", 8}

	return datatype
}

var GlobalDataType = InitDataTypeDef()

type PyPacket struct {
	DstPhyAddr  uint16
	SrcPhyAddr  uint16
	LinkState   byte
	ExpMoreCode byte
	Priority    byte
	HopCnt      byte
	HiProtocode byte

	MessageType  byte
	TranNbr      byte
	SecurityCode uint16
	FileName     string
	CloseFlag    byte
	FileOffset   uint
	Swath        uint16

	SrcNodeId  uint16
	DstNodeId  uint16
	MsgType    byte
	IsRouter   byte
	HopMetric  byte
	VerifyIntv uint16

	RespCode byte

	Raw      []byte
	FileData []byte

	Adjustment []uint
	Time       []uint

	CollectMode byte
	P1          uint
	P2          uint
	TableNbr    uint16
	TableDefSig uint16
}



func (p *PyPacket) Decode_pkt(pkt []byte) {
	//fmt.Println("Decode_Pkt, input pkt size : ", len(pkt))

	rehead := make([]uint16, 4)
	for i := 0; i < 4; i++ {
		data := []byte{pkt[i*2], pkt[i*2+1]}
		rehead[i] = binary.BigEndian.Uint16(data)
	}

	p.LinkState = byte(rehead[0] >> 12)
	p.DstPhyAddr = uint16(rehead[0]) & 0x0FFF
	p.ExpMoreCode = byte((uint16(rehead[1]) & 0xC000) >> 14)
	p.Priority = byte((uint16(rehead[1]) & 0x3000) >> 12)
	p.SrcPhyAddr = uint16(rehead[1]) & 0x0FFF
	p.HiProtocode = byte(rehead[2] >> 12)
	p.DstNodeId = uint16(rehead[2]) & 0x0FFF
	p.HopCnt = byte(rehead[3] >> 12)
	p.SrcNodeId = uint16(rehead[3]) & 0x0FFF

	remsg := pkt[8:]
	p.MsgType = remsg[0]
	p.TranNbr = remsg[1]

	if p.HiProtocode == 0 {
		switch p.MsgType {
		case 0x9:
			fmt.Printf("Message Type is : %x\n", p.MsgType)
			//p.Msg_hello(remsg[2:])
			break
		case 0x89:
			fmt.Printf("what SMessage Type is : %x\n", p.MsgType)
			p.Msg_hello(remsg[2:])
			break
		case 0x8f:
			fmt.Printf("Message Type is : %x\n", p.MsgType)
			break
		case 0x90:
			fmt.Printf("Message Type is : %x\n", p.MsgType)
			break
		case 0x93:
			fmt.Printf("Message Type is : %x\n", p.MsgType)
			break
		default:
			fmt.Printf("Message Type is : %x\n", p.MsgType)
			p.Raw = remsg[2:]
			break
		}
	} else if p.HiProtocode == 1 {
		switch p.MsgType {
		case 0x89:
			fmt.Printf("Message Type is : %x\n", p.MsgType)
			p.Msg_collectdata_response(remsg[2:])
		case 0x9d:
			fmt.Printf("Message Type is : %x\n", p.MsgType)
			p.Msg_fileupload_response(remsg[2:])
			break
		case 0x97:
			fmt.Printf("Message Type is : %x\n", p.MsgType)
			p.Msg_clock_response(remsg[2:])
			break
		default:
			fmt.Printf("Message Type is : %x\n", p.MsgType)
			break

		}
	}

}

func (p *PyPacket) Encode_bin(encode_flag string, size int) ([]byte) {

	switch encode_flag {
	case "":
		return nil
		break
	case "Pkt_fileupload_cmd":
		result := make([]byte, size+1)
		result[0] = 0x1d
		result[1] = p.TranNbr
		result[2], result[3] = Uint16ToByte(p.SecurityCode)
		str_buffer := StringToByte(p.FileName)
		//fmt.Println("String Size : ", len(str_buffer))
		for i := 0; i < len(str_buffer); i++ {
			result[i+4] = str_buffer[i]
		}
		result[len(str_buffer)+4] = '\000'
		result[len(str_buffer)+5] = p.CloseFlag
		result[len(str_buffer)+6], result[len(str_buffer)+7], result[len(str_buffer)+8], result[len(str_buffer)+9] = UintToByte(p.FileOffset)
		result[len(str_buffer)+10], result[len(str_buffer)+11] = Uint16ToByte(p.Swath)
		return result
		break
	case "Pkt_hello_response": //pkt pkt_hello_response encode_bin(['Byte', 'Byte', 'Byte', 'Byte', 'UInt2'], [0x89, TranNbr, IsRouter, HopMetric, VerifyIntv])
		result := make([]byte, size)
		p.HopMetric = 0x02
		p.VerifyIntv = 1800
		//fmt.Println("result size : ", len(result))
		result[0] = 0x89
		result[1] = p.TranNbr
		result[2] = p.IsRouter
		result[3] = p.HopMetric
		result[4], result[5] = Uint16ToByte(p.VerifyIntv)



		return result
		break
	case "Pkt_Clock_cmd":
		result := make([]byte, size)
		result[0] = 0x17
		result[1] = p.TranNbr
		result[2], result[3] = Uint16ToByte(p.SecurityCode)
		result[4], result[5], result[6], result[7] = UintToByte(p.Adjustment[0])
		result[8], result[9], result[10], result[11] = UintToByte(p.Adjustment[1])
		return result
		break
	case "Msg_clock_response":
		result := make([]byte, size)
		result[0] = 0x17
		result[1], result[2], result[3], result[4] = UintToByte(p.Adjustment[0])
		result[5], result[6], result[7], result[8] = UintToByte(p.Adjustment[1])
		return result
		break

	case "Pkt_collectdata_cmd":
		result := make([]byte, size)
		result[0] = 0x09
		result[1] = p.TranNbr
		result[2], result[3] = Uint16ToByte(p.SecurityCode)
		result[4] = 0x05
		result[5], result[6] = Uint16ToByte(p.TableNbr)
		result[7], result[8] = Uint16ToByte(p.TableDefSig)
		result[9], result[10], result[11], result[12] = UintToByte(p.P1)
		result[13], result[14] = Uint16ToByte(0) //fildlist

		return result
		break
	case "Pkt_bye_cmd":
		result := make([]byte, size)
		result[0] = 0x0d
		result[1] = 0x0
		return result
	default:
		return nil
		break

	}
	return nil
}

func (p *PyPacket) Pkt_fileupload_cmd(DstNodeId uint16, SrcNodeId uint16, FileName string, Fileoffset uint, TranNbr byte, CloseFlag byte) ([]byte, byte) {
	var pkt []byte
	p.DstNodeId = DstNodeId
	p.SrcNodeId = SrcNodeId
	p.SecurityCode = 0x0000
	p.FileOffset = Fileoffset
	p.TranNbr = TranNbr
	p.CloseFlag = CloseFlag
	p.Swath = 0x0200
	p.HiProtocode = 0x1
	p.FileName = FileName
	p.RespCode = 0x0e

	//fmt.Println("???FileOffset check : ", p.FileOffset)

	size := 1 + 1 + 2 + len(FileName) + 1 + 4 + 2
	hdr := p.PakBus_hdr(p.DstNodeId, p.SrcNodeId, 0x1)
	msg := p.Encode_bin("Pkt_fileupload_cmd", size)

	pkt = hdr
	for i := 0; i < len(msg); i++ {
		pkt = append(pkt, msg[i])
	}
	a, b := Uint16ToByte(CalcSigNullifier(CalcSigFor(pkt, 0xAAAA)))
	if p.FileOffset == 29184{
		a = 0xdb
	}

	pkt = append(pkt, a, b)

	return pkt, p.TranNbr
}

func (p *PyPacket) Pkt_hello_response(DstNodeId uint16, SrcNodeId uint16, TranNbr byte) ([]byte) {
	fmt.Println("Hello Response")
	var pkt []byte

	hdr := p.PakBus_hdr(DstNodeId, SrcNodeId, 0x0)
	msg := p.Encode_bin("Pkt_hello_response", 1+1+1+1+2)
	//fmt.Println("msg size : ", msg)
	pkt = hdr
	for i := 0; i < len(msg); i++ {
		pkt = append(pkt, msg[i])
	}

	//fmt.Println("pkt : ", pkt)
	a, b := Uint16ToByte(CalcSigNullifier(CalcSigFor(pkt, 0xAAAA)))
	//fmt.Printf("cals : %x, %x\n", a, b)
	pkt = append(pkt, a, b)

	return pkt
}

func (p *PyPacket) PakBus_hdr(DstNodeId uint16, SrcNodeId uint16, HiProtoCode byte) ([]byte) {
	p.LinkState = 0xA
	p.ExpMoreCode = 0x2
	p.Priority = 0x1
	p.HopCnt = 0x0
	p.HiProtocode = HiProtoCode
	p.DstPhyAddr = DstNodeId
	p.SrcPhyAddr = SrcNodeId

	temp := make([]uint16, 4)
	temp[0] = uint16(p.LinkState&0xF)<<12 | uint16(p.DstPhyAddr&0xFFF)
	temp[1] = uint16(p.ExpMoreCode&0x3)<<14 | uint16(p.Priority&0x3)<<12 | uint16(p.SrcPhyAddr&0xFFF)
	temp[2] = uint16(p.HiProtocode&0xF)<<12 | uint16(p.DstPhyAddr&0xFFF)
	temp[3] = uint16(p.HopCnt&0xF)<<12 | uint16(p.SrcPhyAddr&0xFFF)

	buffer := make([]byte, 4*2)
	for i := 0; i < 4; i++ {
		buffer[i*2], buffer[i*2+1] = Uint16ToByte(temp[i])
	}
	return buffer
}

//Message
func (p *PyPacket) Msg_hello(pkt []byte) {
	p.IsRouter = pkt[0]
	p.HopMetric = pkt[1]
	p.VerifyIntv = ByteToUint16(pkt[2], pkt[3])
	p.Raw = pkt[4:]
}

func (p *PyPacket) Msg_devconfig_get_settings_response(pkt []byte) {
}

func (p *PyPacket) Msg_devconfig_set_settings_response(pkt []byte) {
}

func (p *PyPacket) Msg_devconfig_control_response(pkt []byte) {
}

func (p *PyPacket) Msg_collectdata_response(pkt []byte) {
	p.RespCode = pkt[0]
	p.Raw = pkt[1:]
}

func (p *PyPacket) Msg_clock_response(pkt []byte) ([]byte) {
	return p.Encode_bin("Msg_clock_response", 1+8)
}

func (p *PyPacket) Msg_getprogstat_response(pkt []byte) {
}

func (p *PyPacket) Msg_getvalues_response(pkt []byte) {
}

func (p *PyPacket) Msg_filedownload_response(pkt []byte) {
}

func (p *PyPacket) Msg_fileupload_response(pkt []byte) {

	p.RespCode = pkt[0]
	fmt.Println("RespCode : ", p.RespCode)
	p.FileOffset = uint(ByteToUint(pkt[1], pkt[2], pkt[3], pkt[4]))
	if len(pkt) > 5 {
		p.FileData = pkt[5:]
	} else {
		p.FileData = nil
	}

}

func (p *PyPacket) Msg_filecontrol_response(pkt []byte) {
}

func (p *PyPacket) Msg_pleasewait(pkt []byte) {
}

func (p *PyPacket) Print_Packet() {
	fmt.Println("Priority : ", p.Priority)
	fmt.Println("HiProtoCode : ", p.HiProtocode)
	fmt.Println("ExpMoreCode : ", p.ExpMoreCode)
	fmt.Println("SrcNodeId : ", p.SrcNodeId)
	fmt.Println("HopCnt : ", p.HopCnt)
	fmt.Println("DstNodeId : ", p.DstNodeId)
	fmt.Println("SrcPhyAddr : ", p.SrcPhyAddr)
	fmt.Println("DstPhyAddr : ", p.DstPhyAddr)
	fmt.Println("LinkState : ", p.LinkState)
	fmt.Println("MsgType : ", p.MsgType)
	fmt.Println("TranNbr : ", p.TranNbr)
}

func (p *PyPacket) Pkt_Clock_cmd(DstNodeId uint16, SrcNodeId uint16, Adjustment []uint, SecuritiCode uint16, TranNbr byte) ([]byte, int) {

	var pkt []byte
	p.DstNodeId = DstNodeId
	p.SrcNodeId = SrcNodeId
	p.Adjustment = Adjustment
	p.TranNbr = TranNbr
	p.SecurityCode = SecuritiCode

	hdr := p.PakBus_hdr(DstNodeId, SrcNodeId, 0x1)
	msg := p.Encode_bin("Pkt_Clock_cmd", 1+1+2+8)

	pkt = hdr
	for i := 0; i < len(msg); i++ {
		pkt = append(pkt, msg[i])
	}
	return pkt, len(pkt)
}

func (p *PyPacket) Pkt_Bye_Cmd(DstNodeId uint16, SrcNodeId uint16) ([]byte, int){

	hdr := p.PakBus_hdr(DstNodeId, SrcNodeId, p.HiProtocode)
	msg := p.Encode_bin("Pkt_Bye_Cmd", 2)

	pkt := hdr
	for i := 0 ; i < len(msg) ; i++{
		pkt = append(pkt, msg[i])
	}
	return pkt, len(pkt)

}
func (p *PyPacket) Nsec_To_Time() () {

}

func (p *PyPacket) Pkt_collectdata_cmd(DstNodeId uint16, SrcNodeId uint16, TableNbr uint16, TableDefSig uint16) ([]byte, byte) {
	var pkt []byte
	TranNbr := NewTranNbr()
	p.DstNodeId = DstNodeId
	p.SrcNodeId = SrcNodeId
	p.TranNbr = TranNbr
	p.TableNbr = TableNbr
	p.TableDefSig = TableDefSig
	p.CollectMode = 0x5
	p.P1 = 1
	p.P2 = 0
	p.SecurityCode = 0x0000

	hdr := p.PakBus_hdr(DstNodeId, SrcNodeId, 0x1)
	msg := p.Encode_bin("Pkt_collectdata_cmd", 1+1+2+1+2+2+4+2)
	pkt = hdr

	for i := 0; i < len(msg); i++ {
		pkt = append(pkt, msg[i])
	}
	a, b := Uint16ToByte(CalcSigNullifier(CalcSigFor(pkt, 0xAAAA)))
	pkt = append(pkt, a, b)
	return pkt, TranNbr
}

func ByteToUint(a byte, b byte, c byte, d byte) (uint32) {
	data := []byte{a, b, c, d}
	return binary.BigEndian.Uint32(data)
}

func ByteToUint16(a byte, b byte) (uint16) {
	data := []byte{a, b}
	return binary.BigEndian.Uint16(data)
}

func Uint16ToByte(temp uint16) (byte, byte) {
	return byte(temp >> 8), byte(temp & 0xFF)
}

func UintToByte(temp uint) (byte, byte, byte, byte) {
	return byte(temp >> 24), byte(temp >> 16), byte(temp >> 8), byte(temp)
}

func StringToByte(temp string) ([]byte) {
	//fmt.Println("Check String to byte", temp)
	return []byte(temp)
}

func NewTranNbr() (byte) {
	transact += 1
	transact &= 0xFF
	return byte(transact)
}

func CalcSigFor(buff []byte, seed uint16) (uint16) {

	sig := seed
	for i := 0; i < len(buff); i++ {
		x := uint16(buff[i])
		if x == 0 && len(buff) == 2{
			continue
		}
		j := sig
		sig = (sig << 1) & 0x1FF
		if sig >= 0x100 {
			sig += 1
		}
		sig = ((((sig + (j >> 8) + x) & 0xFF) | (j << 8))) & 0xFFFF
	}
	return sig

}

func CalcSigNullifier(sig uint16) (uint16) {

	var nulb uint16
	var nullif uint16
	var sig2 uint16
	b2buff := make([]byte, 2)

	b2buff[0], b2buff[1] = Uint16ToByte(nulb)
	nullifbuffer := make([]byte, 2)

	for i:= 0 ; i < 2 ; i++{
		sig = CalcSigFor(b2buff, sig)
		sig2 = (sig<<1) & 0x1FF
		if sig2 >= 0x100 {
			sig2 += 1
		}
		nulb = ((0x100 - (sig2 + (sig >>8))) & 0xFF)
		b2buff[0], b2buff[1] = Uint16ToByte(nulb)
		nullif += nulb
		nullifbuffer[i] = byte(nulb)
	}
	b2buff[0], b2buff[1] = Uint16ToByte(nullif)
	nullif = binary.BigEndian.Uint16(nullifbuffer)
	return nullif



	/*var new_seed uint16 = uint16((sig << 1) & uint16(0x1FF))
	null1 := make([]byte, 1)
	var new_sig int = int(sig)

	if (new_seed >= 0x0100) {
		new_seed++
	}
	null1[0] = byte(uint16(0x0100 - (new_seed + (sig >> 8))))
	new_sig = int(CalcSigFor(null1, uint16(sig)))

	var null2 uint16

	new_seed = uint16(uint16(new_sig<<1) & uint16(0x01FF))
	if new_seed >= 0x0100 {
		new_seed++
	}
	null2 = uint16(uint16(0x0100 - (new_seed + uint16(new_sig>>8))))

	rtn := uint16(null1[0])
	rtn <<= 8
	rtn += null2
	fmt.Println("nullif : ", rtn)
	return rtn*/
}

func quote(pkt []byte) ([]byte) {
	var str string
	for i := 0; i < len(pkt); i++ {
		str += fmt.Sprintf("%0.2x", pkt[i])
	}
	str = strings.Replace(str, "bc", "bcdc", 100)
	str = strings.Replace(str, "bd", "bcdd", 100)
	result, _ := hex.DecodeString(str)
	return result
}

func uquote(pkt []byte) ([]byte) {
	var str string
	for i := 0; i < len(pkt); i++ {
		str += fmt.Sprintf("%0.2x", pkt[i])
	}
	str = strings.Replace(str, "bcdd", "bd", 100)
	str = strings.Replace(str, "bcdc", "bc", 100)
	result, _ := hex.DecodeString(str)
	return result[:len(result)-2]
}