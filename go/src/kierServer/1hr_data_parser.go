package main

import (
	"time"
	"fmt"
)

type Hour_data struct {
	TIMESTAMP time.Time
	Recode int
	WS_ms_S_WVT float32
	WindDir_D1_WVT float32
	AirTc_Avg float32
	RH float32
	Rain_mm_Tot float32
	HOR_sun_Avg float32
	HOR_shad_Avg float32
	DIRN_Avg float32
}
func (hr *Hour_data) PrintHourData(){
	fmt.Println(hr)
}

func (hr *Hour_data) ParseHourData(temp []string){
}