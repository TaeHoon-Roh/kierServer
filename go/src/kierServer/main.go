package main

import (
	"fmt"
	"os"
	"usepackbus"
	"sync"
	"net"
	"bufio"
	"time"
)

var yy_ip = "106.249.253.227:6785"
var buan_ip = "210.223.199.221:6785"
var seoul_ip = "147.46.138.187:6785"
var jeju_ip = "1.209.192.61:6785"
var deagu_ip = "118.45.205.194:6785"
var Input_File_Name = ".TDF"

type DataSet struct {
	City_Name string
	City_ip   string
	File_Path string
	TableList usepackbus.TableDef
	TableData usepackbus.TableDef
	FileList  usepackbus.FileDef
	FileData  usepackbus.FileDef
}

func main() {

	var wait sync.WaitGroup
	wait.Add(4)

	DataSet := make([]DataSet, 4)

	DataSet[0].City_Name = "YangYang"
	DataSet[1].City_Name = "Buan"
	DataSet[2].City_Name = "Seoul"
	DataSet[3].City_Name = "JeJu"

	DataSet[0].City_ip = "106.249.253.227:6785"
	DataSet[1].City_ip = "210.223.199.221:6785"
	DataSet[2].City_ip = "147.46.138.187:6785"
	DataSet[3].City_ip = "1.209.192.61:6785"

	ConnSet := make([]*net.TCPConn, 4)

	ConnSet[0] = usepackbus.ConnectDevice(DataSet[0].City_ip)
	ConnSet[1] = usepackbus.ConnectDevice(DataSet[1].City_ip)
	ConnSet[2] = usepackbus.ConnectDevice(DataSet[2].City_ip)
	ConnSet[3] = usepackbus.ConnectDevice(DataSet[3].City_ip)

	for i := 0; i < 4; i++ {
		DataSet[i].File_Path = "/home/uxfac/Documents"
	}

	/*
		go func(){
			DataSet[0].TableList = usepackbus.GetList(seoul_ip, Input_File_Name)
		} ()
		go func(){
			DataSet[1].TableList = usepackbus.GetList(yy_ip, Input_File_Name)
		} ()
		go func(){
			DataSet[2].TableList = usepackbus.GetList(buan_ip, Input_File_Name)
		} ()
		go func(){
			DataSet[3].TableList = usepackbus.GetList(jeju_ip, Input_File_Name)
		} ()
	*/
	//testPython(2)
	//DataSet[0].TableList = usepackbus.GetTableList(ConnSet[0], Input_File_Name)
	//DataSet[1].TableList = usepackbus.GetTableList(ConnSet[1], Input_File_Name)
	//DataSet[2].TableList = usepackbus.GetTableList(ConnSet[2], Input_File_Name)
	//DataSet[3].TableList = usepackbus.GetTableList(ConnSet[3], Input_File_Name)
	DataSet[0].FileList = usepackbus.GetFileList(ConnSet[0], ".DIR")
	//DataSet[1].FileList = usepackbus.GetFileList(ConnSet[1], ".DIR")
	//DataSet[2].FileList = usepackbus.GetFileList(ConnSet[2], ".DIR")
	//DataSet[3].FileList = usepackbus.GetFileList(jeju_ip, ".DIR")
	//fmt.Println("Check",DataSet[0].TableList.Table_Count)
	//usepackbus.Collect_Data(ConnSet[0], DataSet[0].TableList, "Public")


	//SaveFile(ConnSet[0], DataSet[0])

}

func SaveFile(cnn *net.TCPConn, set DataSet) {

	cnt := set.FileList.FileDef_Counter
	t_flag := time.Now().AddDate(0, 0, -1)
	for {
		file_last_update := set.FileList.FileDef[cnt-1].LastUpdate
		file_name := set.FileList.FileDef[cnt-1].FileName
		if t_flag.Unix() > file_last_update.Unix() {
			break
		} else {
			fmt.Println(file_last_update)

			_, filepath := MakeFileName(set, file_last_update, file_name)
			fmt.Println(filepath)
			file, err := os.OpenFile(
				filepath,
				os.O_CREATE|os.O_RDWR|os.O_TRUNC, // 파일이 없으면 생성,
				// 읽기/쓰기, 파일을 연 뒤 내용 삭제
				os.FileMode(0644), // 파일 권한은 644
			)

			if err != nil {
				fmt.Println(err)
			}

			FileBuffer := usepackbus.Collect_Data_File(cnn, set.FileList.FileDef[cnt-1].FileName)

			w := bufio.NewWriter(file)
			w.Write(FileBuffer)
			w.Flush()
			file.Close()

			cnt--
		}
	}

}

/*
func Clock_Sync() {
	cnn := tcptest()
	if TranNbr == '\000' {
		TranNbr = usepackbus.NewTranNbr()
	}

	out_packet := usepackbus.PyPacket{}
	in_packet := usepackbus.PyPacket{}
	temp, _ := out_packet.Pkt_Clock_cmd(1, 2050, []uint{0, 0}, 0x0000, TranNbr)
	fmt.Println()

	in_buffer := bufio.NewReader(cnn)
	out_buffer := bufio.NewWriter(cnn)

	for i := 0; i < 1; i++ {
		t1 := time.Now().Nanosecond()
		tcpSendBuffer(out_buffer, temp)
		//retime := time.Time{}
		readbuffer, size := tcpReadBuffer(in_buffer)
		fmt.Println("ReadData !!!!", "Size is : ", size)
		//in_buffer.Reset(cnn)
		if size == 0 {
			fmt.Println("null")
		}
		for i := 0; i < size; i++ {
			fmt.Printf("\\x%0.2x", readbuffer[i])
		}
		fmt.Println()
		in_packet.Decode_pkt(readbuffer)
		t2 := time.Now().Nanosecond()
		delay := (t2 - t1) / 2
		fmt.Println("Delay is : ", delay)
		//longtime := in_packet.Nsec_To_Time()
	}

}
*/

/*
func tcptest() (*net.TCPConn) {
	service := yy_ip

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	//conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	checkError(err)

	return conn
}

*/
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s", err.Error())
		os.Exit(1)
	}

}

/*
func testPython(choice int) {
	str_test := "\xBD\xa0\x01\x98\x02\x10\x01\x08\x02\x09\x02\x00\x00\x05\x00\x08\xe9\xfa\x00\x00\x00\x01\x00\x00\"\xe8\xBD"
	str1 := "\xBD\x90\x01X\x02\x00\x01\x08\x02\t\x01\x00\x02\x07\x08\xf6\x86\xBD"
	str2 := "\xBD\xa0\x01\x98\x02\x10\x01\x08\x02\x1d\x02\x00\x00.TDF\x00\x00\x00\x00\x00\x00\x02\x00\xeet\xBD"
	str3 := "\xBD\xa0\x01\x98\x02\x10\x01\x08\x02\x1d\x02\x00\x00.TDF\x00\x00\x00\x00\x02\x00\x02\x00\xb6]\xBD"
	str4 := "\xBD\xa0\x01\x98\x02\x10\x01\x08\x02\x1d\x02\x00\x00.TDF\x00\x00\x00\x00\x04\x00\x02\x00|E\xBD"
	str5 := "\xBD\xa0\x01\x98\x02\x10\x01\x08\x02\x1d\x02\x00\x00.TDF\x00\x00\x00\x00\x06\x00\x02\x00B-\xBD"
	str6 := "\xBD\xa0\x01\x98\x02\x10\x01\x08\x02\x1d\x02\x00\x00.TDF\x00\x00\x00\x00\x08\x00\x02\x00\x08\x15\xBD"
	str7 := "\xBD\xa0\x01\x98\x02\x10\x01\x08\x02\x1d\x02\x00\x00.TDF\x00\x00\x00\x00\n\x00\x02\x00\xcf\xfd\xBD"
	str8 := "\xBD\xa0\x01\x98\x02\x10\x01\x08\x02\x1d\x02\x00\x00.TDF\x00\x00\x00\x00\x0c\x00\x02\x00\x95\xe5\xBD"
	str9 := "\xBD\xa0\x01\x98\x02\x10\x01\x08\x02\x1d\x02\x00\x00.TDF\x00\x00\x00\x00\r\xa7\x02\x00\xa7\x97\xBD"

	if choice == 2 {
		inpack := usepackbus.PyPacket{}
		temp := testPacket(str_test)
		cnn := tcptest()
		cnn.Write(temp)
		data := make([]byte, 4096)
		size, _ := cnn.Read(data)
		inpack.Decode_pkt(data)
		fmt.Println("\nsize is : ", size)
		fmt.Printf("\n%x\t", data[:size])
		fmt.Println("Message Type", inpack.MsgType)
	} else {
		temp := testPacket(str1)
		cnn := tcptest()
		cnn.Write(temp)
		data := make([]byte, 4096)
		size, _ := cnn.Read(data)
		fmt.Printf("\n%x\t", data[:size])

		temp = testPacket(str2)
		cnn.Write(temp)
		size, _ = cnn.Read(data)
		fmt.Printf("\n%x\t", data[:size])

		temp = testPacket(str3)
		cnn.Write(temp)
		size, _ = cnn.Read(data)
		fmt.Printf("\n%x\t", data[:size])

		temp = testPacket(str4)
		cnn.Write(temp)
		size, _ = cnn.Read(data)
		fmt.Printf("\n%x\t", data[:size])

		temp = testPacket(str5)
		cnn.Write(temp)
		size, _ = cnn.Read(data)
		fmt.Printf("\n%x\t", data[:size])

		temp = testPacket(str6)
		cnn.Write(temp)
		size, _ = cnn.Read(data)
		fmt.Printf("\n%x\t", data[:size])

		temp = testPacket(str7)
		cnn.Write(temp)
		size, _ = cnn.Read(data)
		fmt.Printf("\n%x\t", data[:size])

		temp = testPacket(str8)
		cnn.Write(temp)
		size, _ = cnn.Read(data)
		fmt.Printf("\n%x\t", data[:size])

		temp = testPacket(str9)
		cnn.Write(temp)
		size, _ = cnn.Read(data)
		fmt.Printf("\n%x\t", data[:size])
	}
}

func testPacket(str string) ([]byte) {

	buff := strings.Split(str, "\\")
	fmt.Println("buff size : ", len(buff))
	for i := range buff {
		fmt.Printf("%x\t", buff[i])
	}
	fmt.Println()

	temp := []byte(buff[0])

	for i := range temp {
		fmt.Printf("%x\t", temp[i])
	}

	return temp
}
*/
