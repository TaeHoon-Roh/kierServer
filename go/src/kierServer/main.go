package main

import (
	"fmt"
	"os"
	"usepackbus"
	"net"
	"bufio"
	"time"
	"strconv"
	"strings"
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

/*	cnt := DataSet[0].FileList.FileDef_Counter
	dirpath, filepath := MakeFileName(DataSet[0], DataSet[0].FileList.FileDef[cnt-1].LastUpdate, "CRD:YangYang.1min_data_35.dat")

	os.MkdirAll(dirpath, os.ModePerm)
	fmt.Println("File Path : ",filepath)
	file, err := os.OpenFile(
		filepath,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC, // 파일이 없으면 생성,
		// 읽기/쓰기, 파일을 연 뒤 내용 삭제
		os.FileMode(0644), // 파일 권한은 644
	)

	if err != nil {
		fmt.Println(err)
	}

	FileBuffer := usepackbus.Collect_Data_File(ConnSet[0], "CRD:YangYang.1min_data_35.dat")

	w := bufio.NewWriter(file)
	w.Write(FileBuffer)
	w.Flush()
	file.Close()*/
	FileBuffer := usepackbus.Collect_Data_File(ConnSet[0], "CRD:YangYang.1sec_data_30.dat")
	fmt.Println(FileBuffer)
	//SaveFile(ConnSet[2], DataSet[2])

}

func SaveFile(cnn *net.TCPConn, set DataSet) {

	cnt := set.FileList.FileDef_Counter
	t_flag := time.Now().AddDate(0, 0, -1)
	for {
		file_last_update := set.FileList.FileDef[cnt-1].LastUpdate
		file_name := set.FileList.FileDef[cnt-1].FileName
		file_name = "CRD:Seoul.1sec_data_38.dat"
		if t_flag.Unix() > file_last_update.Unix() {
			break
		} else {
			fmt.Println(file_last_update)

			dirpath, filepath := MakeFileName(set, file_last_update, file_name)
			os.MkdirAll(dirpath, os.ModePerm)
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

			FileBuffer := usepackbus.Collect_Data_File(cnn, file_name)

			w := bufio.NewWriter(file)
			w.Write(FileBuffer)
			w.Flush()
			file.Close()

			cnt--
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s", err.Error())
		os.Exit(1)
	}

}

func MakeFileName(set DataSet, t time.Time, filename string) (string, string) {
	fmt.Println("MakeFileName!!")
	str := "/home/uxfac/Documents/"
	str += set.City_Name + "/"
	year := t.Year()
	fmt.Println("Year : ", year)
	str += strconv.Itoa(year) + "/"
	month := int(t.Month())
	str += strconv.Itoa(month)
	dirpath := str


	filenamebuff := strings.Split(filename,".")
	filenamebuffer := strings.Split(filenamebuff[1], "_")

	str += "/"+filenamebuff[0] + "_" + filenamebuffer[1] + "_"+ t.String() + ".dat"
	filepath := str

	return dirpath, filepath
}