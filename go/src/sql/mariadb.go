package sql

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() {
	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "root:uxfac@tcp(127.0.0.1:3306)/kierserverbackup")
	if err != nil {
		fmt.Println("Db Connect Error!!")
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Db Connect Error!!")
	}
	fmt.Println("Db is : ", db)
	// Connect and check the server version
	CreateTable(db)

	var name string
	err = db.QueryRow("show tables").Scan(&name)
	if err != nil {
		fmt.Println("Db Connect Error!!")
	}
	fmt.Println(name)
}

func CreateTable(db *sql.DB) {
	city := []string{"buan", "daegu", "jeju", "seoul", "yangyang", "kier"}
	term := []string{"hour", "min1", "min10", "sec", "sunhour"}
		q1 := "(RegTimeStamp datetime, Record varchar(20), AirTc decimal(10,3), Humidity decimal(10,3), WindDegree decimal(10,3), WindSpeed decimal(10,3), Rain decimal(10,3), HourSun decimal(10,3), HorShad decimal(10,3), HorDirn decimal(10,3), Etc varchar(25))"

		for i:= 0 ; i < 6 ; i ++{
			for j := 0 ; j < 5 ; j++{
				result, err := db.Exec("CREATE TABLE kierweater_"+city[i]+"_"+term[j]+q1)
				if err!= nil{
					fmt.Println("Table Make Error!!")
				} else{
					fmt.Println(result)
				}

			}
		}

	q2 := "(RegTimeStamp datetime, Sec INTEGER, Min1 INTEGER, Min10 INTEGER, Hour INTEGER, SunHour INTEGER, CheckError decimal(10,3))"
	//result, err := db.Exec("CREATE TABLE kierweater_" + city[0] + "_" + term[0] + q2)
	for i := 0; i < 6; i ++ {
		result, err := db.Exec("CREATE TABLE kierweater_" + city[i] +"Manager"+ q2)
		if err != nil {
			fmt.Println("Table Make Error!!")
		} else {
			fmt.Println(result)

		}
	}

}

func Test() {
	fmt.Println("hi")
}
