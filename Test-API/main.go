package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	session    *mgo.Session
	collection *mgo.Collection
)
var record []string

//const MongoDb details
const (
	Host         = "52.64.154.200:3306"
	AuthUserName = "root"
	AuthPassword = "37ARsYfKLGGrUF5"
	AuthDatabase = "admin"
	Collection   = "CoLLections"
)

type Message struct {
	DisplayMessage string `json:"message"`
}

type Number struct {
	SerialNumber string           `json:"serialNumber" db:"serialnumber"`
	AlarmType    string           `json:"alarmType" db:"alarmtype"`
	CreateDate   time.Time        `json:"createTime"`
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
	_, err = db.Query("INSERT INTO Test(serialnumber,alarmtype,createtime,Pfirstname,Plastname,Pemail,Pphone,Pmobile,Sfirstname,Slastname,Semail,Sphone,Smobile) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", data.SerialNumber, data.AlarmType, time.Now().UTC(), data.Primary.FirstName, data.Primary.LastName, data.Primary.Email, data.Primary.Phone, data.Primary.Mobile, data.Secondary.FirstName, data.Secondary.LastName, data.Secondary.Email, data.Secondary.Phone, data.Secondary.Mobile)
	//insert, err := db.Query("INSERT  into Test VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		BadResponse(w, "Data are Not Uploaded")
	}
	CreateSuccessResponse(w, "Data are Uploded")
	defer db.Close()

	/*mongoDBDialInfo := &mgo.DialInfo{
		Addrs:   []string{"45.76.175.38:27017"},

		Timeout: 60 * time.Second,
	}
	a, err := mgo.DialWithInfo(mongoDBDialInfo)
	session = a
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	session.SetMode(mgo.Monotonic, true)
	collection = session.DB("lcr").C("test")

	err = collection.Insert(data)
	if err != nil {
		BadResponse(w, "Data are Not uploaded")
	}*/
	//CreateSuccessResponse(w, "Data are successfull Uploaded")
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
