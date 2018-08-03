package usepackbus

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"net"
	"time"
)

var yy_ip = "106.249.253.227:6785"
var buan_ip = "210.223.199.221:6785"
var seoul_ip = "147.46.138.187:6785"
var jeju_ip = "1.209.192.61:6785"
var deagu_ip = "118.45.205.194:6785"

var file_offset int
var tabledefloop int
var TableDefStr string
var TranNbr byte
var Send_count = 0
var transact = 0

var TableList TableDef
var FileList FileDef
var FileData string

func GetTableList(cnn *net.TCPConn, FileName string) (TableDef) {

	TranNbr = NewTranNbr()

	out_packet := PyPacket{}
	temp, _ := out_packet.Pkt_fileupload_cmd(1, 2050, FileName, 0x00000000, 1, 0x00)
	fmt.Println()

	in_buffer := bufio.NewReader(cnn)
	out_buffer := bufio.NewWriter(cnn)

	tcpSendBuffer(out_buffer, temp)

	in_packet := PyPacket{}

	FileDataBuffer := make([]byte, 1)
	for {

		readbuffer, size := tcpReadBuffer(in_buffer)

		fmt.Println("ReadData !!!!", "Size is : ", size)
		in_buffer.Reset(cnn)
		if size == 0 {
			break
		}
		for i := 0; i < size; i++ {
			fmt.Printf("\\x%0.2x", readbuffer[i])
		}

		fmt.Println()
		fmt.Println(string(readbuffer))
		in_packet.Decode_pkt(readbuffer[0:size])
		TranNbr = in_packet.TranNbr
		//fmt.Println("=====TranNbr : ", TranNbr)
		if in_packet.RespCode != 0 && in_packet.FileData == nil {
			break
		}
		if out_packet.DstNodeId != in_packet.SrcNodeId || out_packet.SrcNodeId != in_packet.DstNodeId {
			continue

		} else {
			if in_packet.HiProtocode == 0 {
				switch in_packet.MsgType {
				case 0x9:
					fmt.Printf("Message Type is : %x and HiProtocal is : %x\n", in_packet.MsgType, in_packet.HiProtocode)
					//in_packet.Print_Packet()
					tcpSendBuffer(out_buffer, in_packet.Pkt_hello_response(in_packet.SrcNodeId, in_packet.DstNodeId, in_packet.TranNbr))
					break
				case 0x89:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				case 0x8f:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				case 0x90:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				case 0x93:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				default:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				}
			} else if in_packet.HiProtocode == 1 {
				switch in_packet.MsgType {
				case 0x9d:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					fmt.Println("FileOffset Check : ", in_packet.FileOffset)
					in_packet.FileOffset = in_packet.FileOffset + uint(len(in_packet.FileData))
					if in_packet.FileData != nil {
						for i := 0; i < len(in_packet.FileData); i++ {
							FileDataBuffer = append(FileDataBuffer, in_packet.FileData[i])
						}
						temp, _ = out_packet.Pkt_fileupload_cmd(1, 2050, FileName, in_packet.FileOffset, TranNbr, 0x00)
						tcpSendBuffer(out_buffer, temp)
					}
					break
				default:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break

				}

			}
		}

	}

	table := TableDef{}
	table.Parse_TableDef(FileDataBuffer)

	table.Print_TableDef()
	fmt.Println()
	//Collect_Data(cnn, table, "SunHouR")

	//lastBuffer, _ := out_packet.Pkt_Bye_Cmd(1, 2050)
	//tcpSendBuffer(out_buffer, lastBuffer)
	//cnn.Close()

	return table

}

func Collect_Data(cnn *net.TCPConn, tabledef TableDef, tablename string) {
	tablenbr := tabledef.GetTableNbr(tablename)
	out_packet := PyPacket{}
	temp, _ := out_packet.Pkt_collectdata_cmd(1, 2050, tablenbr+1, tabledef.TabelHeader[tablenbr].Table_Sig)
	fmt.Println("tablenbr : ", tablenbr)
	fmt.Println("tabledef.TabelHeader[tablenbr].Table_Sig : ", tabledef.TabelHeader[tablenbr].Table_Sig)

	in_buffer := bufio.NewReader(cnn)
	out_buffer := bufio.NewWriter(cnn)

	tcpSendBuffer(out_buffer, temp)

	in_packet := PyPacket{}

	readbuffer, size := tcpReadBuffer(in_buffer)
	in_buffer.Reset(cnn)
	fmt.Println("ReadData !!!!", "Size is : ", size)
	for i := 0; i < size; i++ {
		fmt.Printf("\\x%0.2x", readbuffer[i])
	}
	fmt.Println()
	in_packet.Decode_pkt(readbuffer[0:size])

	fmt.Println(" what Message Type is : ", in_packet.MsgType)

	fmt.Println("Start CollectData")
	frag := tabledef.Parse_CollectData(in_packet.Raw, 8)
	Print_Frag(frag)
}



func GetFileList(cnn *net.TCPConn, FileName string) (FileDef) {

	TranNbr = NewTranNbr()

	out_packet := PyPacket{}
	temp, _ := out_packet.Pkt_fileupload_cmd(1, 2050, FileName, 0x00000000, TranNbr, 0x00)
	fmt.Println()

	in_buffer := bufio.NewReader(cnn)
	out_buffer := bufio.NewWriter(cnn)

	tcpSendBuffer(out_buffer, temp)

	in_packet := PyPacket{}

	FileDataBuffer := make([]byte, 1)
	for {

		readbuffer, size := tcpReadBuffer(in_buffer)
		fmt.Println("FileDataBuffer Size : ", len(FileDataBuffer))
		in_buffer.Reset(cnn)
		/*fmt.Println("ReadData !!!!", "Size is : ", size)
		//in_buffer.Reset(cnn)
		if size == 0 {
			break
		}
		for i := 0; i < size; i++ {
			fmt.Printf("\\x%0.2x", readbuffer[i])
		}*/
		fmt.Println()
		in_packet.Decode_pkt(readbuffer[0:size])
		TranNbr = in_packet.TranNbr
		//fmt.Println("=====TranNbr : ", TranNbr)
		if in_packet.RespCode != 0 && in_packet.FileData == nil {
			break
		}
		if out_packet.DstNodeId != in_packet.SrcNodeId || out_packet.SrcNodeId != in_packet.DstNodeId {
			continue

		} else {
			if in_packet.HiProtocode == 0 {
				switch in_packet.MsgType {
				case 0x9:
					fmt.Printf("Message Type is : %x and HiProtocal is : %x\n", in_packet.MsgType, in_packet.HiProtocode)
					//in_packet.Print_Packet()
					tcpSendBuffer(out_buffer, in_packet.Pkt_hello_response(in_packet.SrcNodeId, in_packet.DstNodeId, in_packet.TranNbr))
					break
				case 0x89:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				case 0x8f:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				case 0x90:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				case 0x93:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				default:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				}
			} else if in_packet.HiProtocode == 1 {
				switch in_packet.MsgType {
				case 0x9d:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					fmt.Println("FileOffset Check : ", in_packet.FileOffset)
					in_packet.FileOffset = in_packet.FileOffset + uint(len(in_packet.FileData))
					if in_packet.FileData != nil {
						for i := 0; i < len(in_packet.FileData); i++ {
							FileDataBuffer = append(FileDataBuffer, in_packet.FileData[i])
						}
						temp, _ = out_packet.Pkt_fileupload_cmd(1, 2050, FileName, in_packet.FileOffset, TranNbr, 0x00)
						tcpSendBuffer(out_buffer, temp)
					}
					break
				default:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break

				}

			}
		}

	}

	fmt.Println("Check dir", " File buffer Size : ", len(FileDataBuffer))
	file := FileDef{}
	file.Parse_FileDir(FileDataBuffer)

	return file

}

func Collect_Data_File(cnn *net.TCPConn, FileName string) ([]byte) {
	fmt.Println("Collect Data File!!", " File Name is : ", FileName)

	if TranNbr == 0 {
		TranNbr = NewTranNbr()
	}
	TranNbr = 2

	out_packet := PyPacket{}
	temp, _ := out_packet.Pkt_fileupload_cmd(1, 2050, FileName, 0x00000000, TranNbr, 0x00)
	fmt.Println()

	in_buffer := bufio.NewReader(cnn)
	out_buffer := bufio.NewWriter(cnn)

	tcpSendBuffer(out_buffer, temp)

	in_packet := PyPacket{}

	FileDataBuffer := make([]byte, 0)
	for {

		readbuffer, size := tcpReadBuffer(in_buffer)
		in_buffer.Reset(cnn)
		time.Sleep(1*time.Nanosecond)
		if size == 0 {
			break
		}
/*
		fmt.Println("ReadData !!!!", "Size is : ", size)
		for i := 0; i < size; i++ {
			fmt.Printf("\\x%0.2x", readbuffer[i])
		}
		fmt.Println()

*/		in_packet.Decode_pkt(readbuffer[0:size])
		TranNbr = in_packet.TranNbr
		//fmt.Println("=====TranNbr : ", TranNbr)
		if in_packet.RespCode != 0 {
			break
		}
		if out_packet.DstNodeId != in_packet.SrcNodeId || out_packet.SrcNodeId != in_packet.DstNodeId {
			continue

		} else {
			if in_packet.HiProtocode == 0 {
				switch in_packet.MsgType {
				case 0x9:
					fmt.Printf("Message Type is : %x and HiProtocal is : %x\n", in_packet.MsgType, in_packet.HiProtocode)
					//in_packet.Print_Packet()
					tcpSendBuffer(out_buffer, in_packet.Pkt_hello_response(in_packet.SrcNodeId, in_packet.DstNodeId, in_packet.TranNbr))
					break
				case 0x89:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				case 0x8f:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				case 0x90:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				case 0x93:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				default:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break
				}
			} else if in_packet.HiProtocode == 1 {
				switch in_packet.MsgType {
				case 0x9d:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)

					in_packet.FileOffset = in_packet.FileOffset + uint(len(in_packet.FileData))
					fmt.Println("File Offset : ", in_packet.FileOffset)
					if in_packet.FileData != nil {
						for i := 0; i < len(in_packet.FileData); i++ {
							FileDataBuffer = append(FileDataBuffer, in_packet.FileData[i])
						}
						temp, _ = out_packet.Pkt_fileupload_cmd(1, 2050, FileName, in_packet.FileOffset, TranNbr, 0x00)
						tcpSendBuffer(out_buffer, temp)
					}
					//fmt.Println("FileOffset Check : ", in_packet.FileOffset)
					break
				default:
					fmt.Printf("Message Type is : %x\n", in_packet.MsgType)
					break

				}

			}
		}

	}

	fmt.Println(string(FileDataBuffer))

	strbuff := strings.Split(string(FileDataBuffer), ",")

	fmt.Println(strbuff)

	return FileDataBuffer
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s", err.Error())
		os.Exit(1)
	}

}
