package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var record []string

type State struct {
	DisplayString string `json:"displayString"`
}
type DeviceData struct {
	Device   string `json:"device"`
	Online   bool   `json:"online"`
	Firmware string `json:"firmware"`
	State    `json:"sp_regstate"`
	Number   `json:"number"`
}

type Message struct {
	DisplayMessage string `json:"message"`
}

type Number struct {
	SerialNumber string           `json:"serialNumber" db:"serialnumber"`
	AlarmType    string           `json:"alarmType" db:"alarmtype"`
	CreateDate   string           `json:"createTime"`
	Primary      PrimaryContact   `json:"primaryContact" db:"primary"`
	Secondary    SecondaryContact `json:"secondaryContact" db"secondary"`
}

type PrimaryContact struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Mobile    string `json:"mobile"`
}
type SecondaryContact struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Mobile    string `json:"mobile"`
}

const (
	host     = "52.64.154.200"
	port     = 3306
	user     = "root"
	password = "lmkt@ptcl"
	dbname   = "mysql"
)

var db *sql.DB

func main() {
	http.HandleFunc("/dataupload", postdata)
	http.HandleFunc("/getdata", getdata)
	http.ListenAndServe(":8080", nil)
}
func postdata(w http.ResponseWriter, req *http.Request) {
	var data Number
	res := req.Body
	fmt.Println(res)
	rep, err := ioutil.ReadAll(res)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(rep, &data)
	db, err := sql.Open("mysql", "root:lmkt@ptcl@tcp(52.64.154.200:3306)/mysql")
	if err != nil {
		fmt.Println(err)
	}
	var id = 0
	_, err = db.Query("INSERT INTO Test(id,serialnumber,alarmtype,createtime,Pfirstname,Plastname,Pemail,Pphone,Pmobile,Sfirstname,Slastname,Semail,Sphone,Smobile) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", id, data.SerialNumber, data.AlarmType, time.Now().UTC(), data.Primary.FirstName, data.Primary.LastName, data.Primary.Email, data.Primary.Phone, data.Primary.Mobile, data.Secondary.FirstName, data.Secondary.LastName, data.Secondary.Email, data.Secondary.Phone, data.Secondary.Mobile)
	if err != nil {
		BadResponse(w, "Data are Not Uploaded")
	} else {
		CreateSuccessResponse(w, "Data are Uploded")
	}
	defer db.Close()

}
func getdata(w http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	id := queryValues.Get("id")
	uptime := queryValues.Get("up")
	totime := queryValues.Get("to")
	db, err := sql.Open("mysql", "root:lmkt@ptcl@tcp(52.64.154.200:3306)/mysql")
	if err != nil {
		fmt.Println(err)
	}
	data, err := db.Query("select * from Devices d, Test t where d.device = t.serialnumber and d.device = ? and t.createtime BETWEEN ? AND ?", id, uptime, totime)
	if err != nil {
		fmt.Print(err)
	}

	emp := DeviceData{}
	res := []DeviceData{}
	for data.Next() {
		var id int
		var online bool
		var device, firmware, createtime, serialnumber, spRegStaee, alarmtype, Pfirstname, Plastname, Pemail, Pphone, Pmobile, Sfirstname, Slastname, Semail, Sphone, Smobile string
		err = data.Scan(&id, &device, &online, &firmware, &spRegStaee, &id, &serialnumber, &alarmtype, &createtime, &Pfirstname, &Plastname, &Pemail, &Pphone, &Pmobile, &Sfirstname, &Slastname, &Semail, &Sphone, &Smobile)
		if err != nil {
			fmt.Println(err)
		}
		emp.Device = device
		emp.Online = online
		emp.Firmware = firmware
		emp.DisplayString = spRegStaee
		emp.SerialNumber = serialnumber
		emp.AlarmType = alarmtype
		emp.CreateDate = createtime
		emp.Primary.FirstName = Pfirstname
		emp.Primary.LastName = Plastname
		emp.Primary.Email = Pemail
		emp.Primary.Phone = Pphone
		emp.Primary.Mobile = Pmobile
		emp.Secondary.FirstName = Sfirstname
		emp.Secondary.LastName = Slastname
		emp.Secondary.Email = Semail
		emp.Secondary.Phone = Sphone
		emp.Secondary.Mobile = Smobile
		res = append(res, emp)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Print(res)
	messagee, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(messagee)

}
func CreateSuccessResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	status := Message{message}
	messagee, err := json.Marshal(status)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(messagee)
}
func BadResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	status := Message{message}
	messagee, err := json.Marshal(status)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(messagee)
}
