package main

import (
	"time"
	"strings"
	"fmt"
	"io/ioutil"
	"strconv"
)

//index 34 ~

type Weather_Data struct{
	LastUpdate time.Time
	Hour_data_count int
	Min10_Data_count int
	Min1_Data_count int
	Sec_Data_count int
	Hour_data_list []Hour_data
	Min10_Data_list []Min10_Data
	Min1_Data_list []Min1_Data
	Sec_Data_list []Sec_Data
}

type Hour_data struct {
	TIMESTAMP time.Time
	Recode string
	WS_ms_S_WVT float32
	WindDir_D1_WVT float32
	AirTc_Avg float32
	RH float32
	Rain_mm_Tot float32
	HOR_sun_Avg float32
	HOR_shad_Avg float32
	DIRN_Avg float32
}

type Min10_Data struct{
	TIMESTAMP time.Time
	Recode string
	WS_ms_S_WVT float32
	WindDir_D1_WVT float32
	AirTc_Avg float32
	RH float32
	Rain_mm_Tot float32
	HOR_sun_Avg float32
	HOR_shad_Avg float32
	DIRN_Avg float32
}

type Min1_Data struct{
	TIMESTAMP time.Time
	Recode string
	WS_ms_S_WVT float32
	WindDir_D1_WVT float32
	AirTc_Avg float32
	RH float32
	Rain_mm_Tot float32
	HOR_sun_Avg float32
	HOR_shad_Avg float32
	DIRN_Avg float32
}

type Sec_Data struct{
	TIMESTAMP time.Time
	Recode string
	WS_ms_S_WVT float32
	WindDir_D1_WVT float32
	AirTc_Avg float32
	RH float32
	Rain_mm_Tot float32
	HOR_sun_Avg float32
	HOR_shad_Avg float32
	DIRN_Avg float32
}

func TimeDataParser(filename string) (Weather_Data){
	var weatherData Weather_Data


	return weatherData
}

func InsertHourData(filename string) ([]Hour_data, int){
	hour_buff := make([]Hour_data,0)
	hour := Hour_data{}
	cnt := 0
	buff, err := ioutil.ReadFile(filename)
	if err != nil{
		fmt.Println("FileOpen Error : ", filename)
	}
	str := strings.Split(string(buff), ",")
	str_buff := make([]string, 0)
	for i := 0 ; i < len(str) ; i++{
		if strings.Contains(str[i], "\n"){
			buff := strings.Split(str[i],"\n")
			for j := 0 ; j < len(buff) ; j++{
				str_buff = append(str_buff, buff[j])
			}
		} else{
			str_buff = append(str_buff, str[i])
		}
	}
	for i := 38 ; i < len(str_buff) ; i++{
		if str_buff[i] == ""{
			break
		}

		hour.TIMESTAMP = RemakeTimeStemp(str_buff[i])
		i++
		hour.Recode = str_buff[i]
		i++
		var temp float64
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		hour.WS_ms_S_WVT = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		hour.WindDir_D1_WVT = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		hour.AirTc_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		hour.RH = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		hour.Rain_mm_Tot = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		hour.HOR_sun_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		hour.HOR_shad_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		hour.DIRN_Avg = float32(temp)
		cnt++
		hour_buff = append(hour_buff, hour)
	}

	return hour_buff, cnt
}

func InsertMin10Data(filename string) ([]Min10_Data, int){
	min10_buff := make([]Min10_Data,0)
	min10 := Min10_Data{}
	cnt := 0
	buff, err := ioutil.ReadFile(filename)
	if err != nil{
		fmt.Println("FileOpen Error : ", filename)
	}
	str := strings.Split(string(buff), ",")
	str_buff := make([]string, 0)
	for i := 0 ; i < len(str) ; i++{
		if strings.Contains(str[i], "\n"){
			buff := strings.Split(str[i],"\n")
			for j := 0 ; j < len(buff) ; j++{
				str_buff = append(str_buff, buff[j])
			}
		} else{
			str_buff = append(str_buff, str[i])
		}
	}
	for i := 38 ; i < len(str_buff) ; i++{
		if str_buff[i] == ""{
			break
		}

		min10.TIMESTAMP = RemakeTimeStemp(str_buff[i])
		i++
		min10.Recode = str_buff[i]
		i++
		var temp float64
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min10.WS_ms_S_WVT = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min10.WindDir_D1_WVT = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min10.AirTc_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min10.RH = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min10.Rain_mm_Tot = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min10.HOR_sun_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min10.HOR_shad_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min10.DIRN_Avg = float32(temp)
		cnt++
		min10_buff = append(min10_buff, min10)
	}

	return min10_buff, cnt
}

func InsertMin1Data(filename string) ([]Min1_Data, int){
	min1_buff := make([]Min1_Data,0)
	min1 := Min1_Data{}
	cnt := 0
	buff, err := ioutil.ReadFile(filename)
	if err != nil{
		fmt.Println("FileOpen Error : ", filename)
	}
	str := strings.Split(string(buff), ",")
	str_buff := make([]string, 0)
	for i := 0 ; i < len(str) ; i++{
		if strings.Contains(str[i], "\n"){
			buff := strings.Split(str[i],"\n")
			for j := 0 ; j < len(buff) ; j++{
				str_buff = append(str_buff, buff[j])
			}
		} else{
			str_buff = append(str_buff, str[i])
		}
	}
	for i := 38 ; i < len(str_buff) ; i++{
		if str_buff[i] == ""{
			break
		}

		min1.TIMESTAMP = RemakeTimeStemp(str_buff[i])
		i++
		min1.Recode = str_buff[i]
		i++
		var temp float64
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min1.WS_ms_S_WVT = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min1.WindDir_D1_WVT = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min1.AirTc_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min1.RH = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min1.Rain_mm_Tot = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min1.HOR_sun_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min1.HOR_shad_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		min1.DIRN_Avg = float32(temp)
		cnt++
		min1_buff = append(min1_buff, min1)
	}

	return min1_buff, cnt
}

func InsertSecData(filename string) ([]Sec_Data, int){
	sec_buff := make([]Sec_Data,0)
	sec := Sec_Data{}
	cnt := 0
	buff, err := ioutil.ReadFile(filename)
	if err != nil{
		fmt.Println("FileOpen Error : ", filename)
	}
	str := strings.Split(string(buff), ",")
	str_buff := make([]string, 0)
	for i := 0 ; i < len(str) ; i++{
		if strings.Contains(str[i], "\n"){
			buff := strings.Split(str[i],"\n")
			for j := 0 ; j < len(buff) ; j++{
				str_buff = append(str_buff, buff[j])
			}
		} else{
			str_buff = append(str_buff, str[i])
		}
	}
	for i := 38 ; i < len(str_buff) ; i++{
		if str_buff[i] == ""{
			break
		}
		sec.TIMESTAMP = RemakeTimeStemp(str_buff[i])
		i++
		sec.Recode = str_buff[i]
		i++
		var temp float64
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		sec.WS_ms_S_WVT = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		sec.WindDir_D1_WVT = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		sec.AirTc_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		sec.RH = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		sec.Rain_mm_Tot = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		sec.HOR_sun_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		sec.HOR_shad_Avg = float32(temp)
		i++
		temp ,_ = strconv.ParseFloat(str_buff[i],32)
		sec.DIRN_Avg = float32(temp)
		cnt++
		sec_buff = append(sec_buff, sec)

	}

	return sec_buff, cnt
}

func RemakeTimeStemp(t string) (time.Time){
	buff := []byte(t)
	str := strings.Replace(t, " ", "T", 1)
	str = strings.Replace(str, "\"", "", 2)
	str = str + "Z"
	t_buff,err := time.Parse(time.RFC3339, str)
	if err != nil{
		fmt.Println("Time Error", t, "buffer size : ", len(buff) , " buff : " , buff)
	}
	return t_buff
}