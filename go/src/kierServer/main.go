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
	"github.com/gin-gonic/gin"
	"net/http"
)

var yy_ip = "106.249.253.227:6785"
var buan_ip = "210.223.199.221:6785"
var seoul_ip = "147.46.138.187:6785"
var jeju_ip = "1.209.192.61:6785"
var deagu_ip = "118.45.205.194:6785"
var Input_File_Name = ".TDF"

type DataSet struct {
	City_Name   string
	City_ip     string
	File_Path   string
	TableList   usepackbus.TableDef
	TableData   usepackbus.TableDef
	FileList    usepackbus.FileDef
	FileData    usepackbus.FileDef
	PublicData  string
	HourData    string
	FileListStr string
}

type Login struct {
	User string `json:"user" binding:"required"`
}

func WebServer(set []DataSet) {
	r := gin.Default()
	r.Static("/", "/home/uxfac/go/templates/")

	r.POST("", func(c *gin.Context) {
		var test Login
		str, _ := c.GetRawData()
		str_c := string(str)
		fmt.Println("Check", str_c)
		c.ShouldBindJSON(test)
		fmt.Println("Test Print : ", test)
		c.JSON(http.StatusOK, gin.H{
			"test": "11111111111111111111111111111111",
		})

	})

	r.POST("/Status", func(c *gin.Context) {
		var test Login
		str, _ := c.GetRawData()
		str_c := string(str)
		fmt.Println(str_c)
		c.ShouldBindJSON(test)
		fmt.Println(test)
		fmt.Println("Post Check", set[0].HourData)

		c.JSON(http.StatusOK, gin.H{
			"city1": set[0].City_Name,
			"raw1":  set[0].PublicData,
			"city2": set[1].City_Name,
			"raw2":  set[1].PublicData,
			"city3": set[2].City_Name,
			"raw3":  set[2].PublicData,
			"city4": set[3].City_Name,
			"raw4":  set[3].PublicData,
		})
	})

	r.POST("/Now", func(c *gin.Context) {
		var test Login
		str, _ := c.GetRawData()
		str_c := string(str)
		fmt.Println(str_c)
		c.ShouldBindJSON(test)
		fmt.Println(test)
		c.JSON(http.StatusOK, gin.H{
			"city1": set[0].City_Name,
			"raw1":  set[0].HourData,
			"city2": set[1].City_Name,
			"raw2":  set[1].HourData,
			"city3": set[2].City_Name,
			"raw3":  set[2].HourData,
			"city4": set[3].City_Name,
			"raw4":  set[3].HourData,
		})
	})

	r.POST("/Database", func(c *gin.Context) {
		var test Login
		str, _ := c.GetRawData()
		str_c := string(str)
		fmt.Println(str_c)
		c.ShouldBindJSON(test)
		fmt.Println(test)
		c.JSON(http.StatusOK, gin.H{
			"city1": set[0].City_Name,
			"raw1":  set[0].FileListStr,
			"city2": set[1].City_Name,
			"raw2":  set[1].FileListStr,
			"city3": set[2].City_Name,
			"raw3":  set[2].FileListStr,
			"city4": set[3].City_Name,
			"raw4":  set[3].FileListStr,
		})
	})

	r.Run(":5000") // listen and serve on 0.0.0.0:8080
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

	DirList()
	filepath := "/home/uxfac/rawData/YangYang/2018/8/1hr/YangYang.1hr.2018-8-6.dat"
	wd := Weather_Data{}
	wd.Hour_data_list,wd.Hour_data_count = InsertHourData(filepath)
	fmt.Println(wd.Sec_Data_count)
	db := ConnectDb("server")
	//CreateTable(db)
	/*tt, err := db.Exec("INSERT INTO t2 (Record1, Record) VALUES ('hihi', 'hihi');")

	if err != nil{
		fmt.Println("Error", err)
	}
	fmt.Println(tt)
*/



	InsertDBTable(db, filepath, wd)

}

func DirList(){
}
func GetTableDef(set *DataSet, cnn *net.TCPConn) {
	set.TableList = usepackbus.GetTableList(cnn, ".TDF", 1)
	set.PublicData = usepackbus.Collect_Data(cnn, set.TableList, "Public")
	set.HourData = usepackbus.Collect_Data(cnn, set.TableList, "hour")
}

func GetFileDef(set *DataSet, cnn *net.TCPConn) {
	if set.City_Name == "JeJu" {
		set.FileListStr = "Error"
	} else {
		set.FileList = usepackbus.GetFileList(cnn, ".DIR", 1)
		set.FileListStr = fmt.Sprintf("%v", set.FileList)
	}
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

	filenamebuff := strings.Split(filename, ".")
	filenamebuffer := strings.Split(filenamebuff[1], "_")

	str += "/" + filenamebuff[0] + "_" + filenamebuffer[1] + "_" + t.String() + ".dat"
	filepath := str

	return dirpath, filepath
}
