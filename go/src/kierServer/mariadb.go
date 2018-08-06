package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"strconv"
	"time"
)

func ConnectDb(flag string) (*sql.DB) {
	// Create the database handle, confirm driver is present

	if flag == "server" {
		db, err := sql.Open("mysql", "root:uxfac@tcp(127.0.0.1:3306)/kierserver")
		if err != nil {
			fmt.Println("Db Connect Error!!")
		}
		if err := db.Ping(); err != nil {
			fmt.Println("Db Connect Error!!")
		}
		fmt.Println("Db is : ", db)
		return db
	} else if flag == "check" {
		db, err := sql.Open("mysql", "root:uxfac@tcp(127.0.0.1:3306)/kierserverCheck")
		if err != nil {
			fmt.Println("Db Connect Error!!")
		}
		if err := db.Ping(); err != nil {
			fmt.Println("Db Connect Error!!")
		}
		fmt.Println("Db is : ", db)
		return db
	} else if flag == "test"{
		db, err := sql.Open("mysql", "root:uxfac@tcp(127.0.0.1:3306)/test")
		if err != nil {
			fmt.Println("Db Connect Error!!")
		}
		if err := db.Ping(); err != nil {
			fmt.Println("Db Connect Error!!")
		}
		fmt.Println("Db is : ", db)
		return db
	}
	return nil

}

func CreateDatabase() {

}

func CreateTable(db *sql.DB) {
	city := []string{"Buan", "Daegu", "JeJu", "Seoul", "YangYang", "Kier"}
	term := []string{"1hr", "1min", "10min", "1sec", "sunhour"}
	q1 := "(IDX bigint(20) NOT NULL auto_increment, RegTimeStamp datetime, Record varchar(20), WS_ms_S_WVT decimal(10,3), WindDir_D1_WVT decimal(10,3), AirTc_Avg decimal(10,3), RH decimal(10,3), Rain_mm_Tot decimal(10,3), HOR_sun_Avg decimal(10,3), HOR_shad_Avg decimal(10,3), DIRN_Avg decimal(10,3), Etc varchar(25), PRIMARY KEY(IDX))"

	for i := 0; i < 6; i ++ {
		for j := 0; j < 5; j++ {
			result, err := db.Exec("CREATE TABLE kierweater_" + city[i] + "_" + term[j] + q1)
			if err != nil {
				fmt.Println("Table Make Error!!", err)
			} else {
				fmt.Println(result)
			}

		}
	}

	/*q2 := "(RegTimeStamp datetime, Sec INTEGER, Min1 INTEGER, Min10 INTEGER, Hour INTEGER, SunHour INTEGER, CheckError decimal(10,3))"
	//result, err := db.Exec("CREATE TABLE kierweater_" + city[0] + "_" + term[0] + q2)
	for i := 0; i < 6; i ++ {
		result, err := db.Exec("CREATE TABLE kierweater_" + city[i] +"Manager"+ q2)
		if err != nil {
			fmt.Println("Table Make Error!!")
		} else {
			fmt.Println(result)

		}
	}*/

}

func InsertDBTable(db *sql.DB, filepath string, data Weather_Data) {
	city := []string{"Buan", "Daegu", "JeJu", "Seoul", "YangYang", "Kier"}
	term := []string{"1hr", "1min", "10min", "1sec", "sunhour"}

	var flag_city string
	var flag_term string

	for i := 0; i < len(city); i++ {
		if strings.Contains(filepath, city[i]) {
			flag_city = city[i]
			break
		}
	}
	for i := 0; i < len(term); i++ {
		if strings.Contains(filepath, term[i]) {
			flag_term = term[i]
			break
		}
	}
	//fmt.Println(data)
	fmt.Println("City : ", flag_city, " Term : ", flag_term)
	fmt.Println("Data set Count", data.Hour_data_count)

	query1 := "INSERT INTO kierweater_" + flag_city + "_" + flag_term + " "
	query3 := "(IDX, RegTimeStamp, Record, AirTc, Humidity, WindDegree, WindSpeed, Rain, HourSun , HorShad, HorDirn, Etc) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"

	query := query1 + query3

	for i := 0; i < data.Sec_Data_count; i++ {
		etc := strconv.FormatFloat(float64(data.Sec_Data_list[i].DIRN_Avg), 'f', 1, 64, )
		result, err := db.Exec(query,
			1,
			data.Sec_Data_list[i].TIMESTAMP,
			data.Sec_Data_list[i].Recode,
			data.Sec_Data_list[i].WS_ms_S_WVT,
			data.Sec_Data_list[i].WindDir_D1_WVT,
			data.Sec_Data_list[i].AirTc_Avg,
			data.Sec_Data_list[i].RH,
			data.Sec_Data_list[i].Rain_mm_Tot,
			data.Sec_Data_list[i].HOR_sun_Avg,
			data.Sec_Data_list[i].HOR_shad_Avg,
			etc, )
		if err != nil {
			fmt.Println("Insert Data Error ", err)
		} else {
			fmt.Println(result)
		}
	}

	for i := 0; i < data.Min1_Data_count; i++ {
		etc := strconv.FormatFloat(float64(data.Min1_Data_list[i].DIRN_Avg), 'f', 1, 64, )
		result, err := db.Exec(query,
			data.Min1_Data_list[i].TIMESTAMP,
			data.Min1_Data_list[i].Recode,
			data.Min1_Data_list[i].WS_ms_S_WVT,
			data.Min1_Data_list[i].WindDir_D1_WVT,
			data.Min1_Data_list[i].AirTc_Avg,
			data.Min1_Data_list[i].RH,
			data.Min1_Data_list[i].Rain_mm_Tot,
			data.Min1_Data_list[i].HOR_sun_Avg,
			data.Min1_Data_list[i].HOR_shad_Avg,
			etc, )
		if err != nil {
			fmt.Println("Insert Data Error ", err)
		} else {
			fmt.Println(result)
		}
	}

	for i := 0; i < data.Min10_Data_count; i++ {
		etc := strconv.FormatFloat(float64(data.Min10_Data_list[i].DIRN_Avg), 'f', 1, 64, )
		result, err := db.Exec(query,
			data.Min10_Data_list[i].TIMESTAMP,
			data.Min10_Data_list[i].Recode,
			data.Min10_Data_list[i].WS_ms_S_WVT,
			data.Min10_Data_list[i].WindDir_D1_WVT,
			data.Min10_Data_list[i].AirTc_Avg,
			data.Min10_Data_list[i].RH,
			data.Min10_Data_list[i].Rain_mm_Tot,
			data.Min10_Data_list[i].HOR_sun_Avg,
			data.Min10_Data_list[i].HOR_shad_Avg,
			etc, )
		if err != nil {
			fmt.Println("Insert Data Error ", err)
		} else {
			fmt.Println(result)
		}
	}

	for i := 0; i < data.Hour_data_count; i++ {
		//etc := strconv.FormatFloat(float64(data.Hour_data_list[i].DIRN_Avg), 'f',1,64,)
		str := TimestempParsing(data.Hour_data_list[i].TIMESTAMP)
		fmt.Println(str)
		result, err := db.Exec(query,
			0,
			str,
			data.Hour_data_list[i].Recode,
			data.Hour_data_list[i].WS_ms_S_WVT,
			data.Hour_data_list[i].WindDir_D1_WVT,
			data.Hour_data_list[i].AirTc_Avg,
			data.Hour_data_list[i].RH,
			data.Hour_data_list[i].Rain_mm_Tot,
			data.Hour_data_list[i].HOR_sun_Avg,
			data.Hour_data_list[i].HOR_shad_Avg,
			data.Hour_data_list[i].DIRN_Avg,
			"", )
		if err != nil {
			fmt.Println("Insert Data Error ", err)
		} else {
			fmt.Println(result)
		}
	}

}
func TimestempParsing(t time.Time) (string) {
	str := t.String()
	str = strings.Replace(str, "T", "", 1)
	str = strings.Replace(str, "Z", "", 1)
	val := []rune(str)
	substring := string(val[:19])
	fmt.Println("STr : ", substring)
	return substring
}

func Test() {
	fmt.Println("hi")
}
