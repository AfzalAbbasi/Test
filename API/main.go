package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonlvhit/gocron"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"
<<<<<<< HEAD
=======
	"time"
>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51
)

var limiter = rate.NewLimiter(2, 5)

type State struct {
	DisplayString string `json:"displayString"`
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
type StatusData struct {
	Data Status `json:"data"`
}
type Status struct {
	Device      string  `json:"device"`
	Online      bool    `json:"online"`
	Firmware    string  `json:"firmware"`
	SpRegStates []State `json:"spRegStates"`
}
type Person struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type TokenData struct {
	Data SessionData `json:"data"`
}
type DeviceData struct {
	Device   string `json:"device"`
	Online   bool   `json:"online"`
	Firmware string `json:"firmware"`
	State    `json:"sp_regstate"`
	Number   `json:"number"`
}
type DeviceDataa struct {
	DeviceData Device `json:"data"`
}
type SessionData struct {
	Token     string `json:"token"`
	Agent     string `json:"agent"`
	Org       string `json:"org"`
	CreateAt  string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
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
type Device []struct {
	ID           string `json:"id" bson:"id"`
	ObiNumber    string `json:"obiNumber" bson:"obiNumber"`
	MacAddress   string `json:"macAddress" bson:"macAddress"`
	SerialNumber string `json:"serialNumber" bson:"serialNumber"`
	Org          string `json:"org"  bson:"org"`
	DeviceType   string `json:"deviceType" bson:"deviceType"`
}

<<<<<<< HEAD
var sql_con, err = sql.Open("")
=======
var i, err = sql.Open()
>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51

/*func Connerction() *sql.DB {
	db, err :=
	if err != nil {
		fmt.Println(err)
	}
	return db
}*/
func main() {
<<<<<<< HEAD
	sql_con.SetMaxOpenConns(-1)
	http.HandleFunc("/dataupload", postdata)
	http.HandleFunc("/getdata", getdata)
	go func() {
		gocron.Every(1).Minute().Do(Devices)
		<-gocron.Start()
	}()
	http.ListenAndServe(":8080", nil)

=======
	i.SetMaxOpenConns(-1)
	http.HandleFunc("/dataupload", postdata)
	http.HandleFunc("/getdata", getdata)
	//http.ListenAndServe(":8080", nil)
	gocron.Every(1).Minute().Do(Devices)
	<-gocron.Start()
>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51
}
func Devices() {
	token := Session()
	value := fmt.Sprintf("%s %s", "Bearer", token)
	client := &http.Client{}
	request, _ := http.NewRequest("GET", "https://api.obitalk.com/api/v1/devices", nil)
	request.Header.Set("authorization", value)
	response, err := client.Do(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	req1, err := ioutil.ReadAll(response.Body)
	var devicedata DeviceDataa
	json.Unmarshal(req1, &devicedata)
	for _, item := range devicedata.DeviceData {
		a := item.ID
		api := fmt.Sprintf("%s/%s/%s", "https://api.obitalk.com/api/v1/devices", a, "quick-values")
		request1, _ := http.NewRequest("POST", api, nil)
		request1.Header.Set("authorization", value)
		response1, err := client.Do(request1)

		if err != nil {
			fmt.Print(err.Error())
		}
		var data StatusData
		req, err := ioutil.ReadAll(response1.Body)
		json.Unmarshal(req, &data)
		if err != nil {
			fmt.Println(err)
		}
		var b []State
		for _, item := range data.Data.SpRegStates {
			a := item.DisplayString
			if a != "" {
				if a != "Service Not Configured" && a != "null" {
					b = append(b, State{DisplayString: a})
				}
			}
		}
		var id = 0
		data.Data.SpRegStates = b
		j := fmt.Sprintf("%s", b)
<<<<<<< HEAD
		n, err := sql_con.Query("INSERT INTO DeviceInformation(device,online,firmware,spRegState,id) VALUES (?,?,?,?,?)", data.Data.Device, data.Data.Online, data.Data.Firmware, j, id)
		 n.Close()
		if err != nil {
			log.Print(err)
		}


	}

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
	if err != nil {

		fmt.Println(err)
	}
	var id = 0
	if data.SerialNumber==""{
		BadResponse(w,"Serial Number field is Required")
		return
	}
	n, err := sql_con.Query("INSERT INTO UserInformation(id,serialnumber,alarmtype,createtime,Pfirstname,Plastname,Pemail,Pphone,Pmobile,Sfirstname,Slastname,Semail,Sphone,Smobile) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", id, data.SerialNumber, data.AlarmType, time.Now().UTC(), data.Primary.FirstName, data.Primary.LastName, data.Primary.Email, data.Primary.Phone, data.Primary.Mobile, data.Secondary.FirstName, data.Secondary.LastName, data.Secondary.Email, data.Secondary.Phone, data.Secondary.Mobile)

	if err != nil {
		BadResponse(w, "Unable to Uplaod Data")
		return
	} else {
		CreateSuccessResponse(w, "Data  Uploded")
	}
	defer n.Close()

=======
		_, err = i.Query("INSERT INTO DeviceInformation(device,online,firmware,spRegState,id) VALUES (?,?,?,?,?)", data.Data.Device, data.Data.Online, data.Data.Firmware, j, id)
		if err != nil {
			log.Print(err)
		}
		defer i.Close()
	}

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
	if err != nil {
		fmt.Println(err)
	}
	var id = 0
	_, err = i.Query("INSERT INTO UserInformation(id,serialnumber,alarmtype,createtime,Pfirstname,Plastname,Pemail,Pphone,Pmobile,Sfirstname,Slastname,Semail,Sphone,Smobile) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", id, data.SerialNumber, data.AlarmType, time.Now().UTC(), data.Primary.FirstName, data.Primary.LastName, data.Primary.Email, data.Primary.Phone, data.Primary.Mobile, data.Secondary.FirstName, data.Secondary.LastName, data.Secondary.Email, data.Secondary.Phone, data.Secondary.Mobile)
	if err != nil {
		BadResponse(w, "Data are Not Uploaded")
	} else {
		CreateSuccessResponse(w, "Data are Uploded")
	}
	defer i.Close()

>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51
}
func getdata(w http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	id := queryValues.Get("id")
	uptime := queryValues.Get("from")
	totime := queryValues.Get("to")
<<<<<<< HEAD
	if (id=="")||(uptime=="")||(totime=="") {
		BadResponse(w,"Enter the Correct URL")
	}
	data, err := sql_con.Query("select * from DeviceInformation d, UserInformation t where d.device = t.serialnumber and d.device = ? and t.createtime BETWEEN ? AND ?", id, uptime, totime)
	if err != nil {
		BadResponse(w, "Unable to Uplaod Data")
		return
=======

	data, err := i.Query("select * from DeviceInformation d, UserInformation t where d.device = t.serialnumber and d.device = ? and t.createtime BETWEEN ? AND ?", id, uptime, totime)
	if err != nil {
		fmt.Print(err)
>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51
	}
	emp := DeviceData{}
	res := []DeviceData{}
	for data.Next() {
<<<<<<< HEAD
		var idd int
		var online bool
		var device, firmware, createtime, serialnumber, spRegStaee, alarmtype, Pfirstname, Plastname, Pemail, Pphone, Pmobile, Sfirstname, Slastname, Semail, Sphone, Smobile string
		err = data.Scan(&idd, &device, &online, &firmware, &spRegStaee, &idd, &serialnumber, &alarmtype, &createtime, &Pfirstname, &Plastname, &Pemail, &Pphone, &Pmobile, &Sfirstname, &Slastname, &Semail, &Sphone, &Smobile)
=======
		var id int
		var online bool
		var device, firmware, createtime, serialnumber, spRegStaee, alarmtype, Pfirstname, Plastname, Pemail, Pphone, Pmobile, Sfirstname, Slastname, Semail, Sphone, Smobile string
		err = data.Scan(&id, &device, &online, &firmware, &spRegStaee, &id, &serialnumber, &alarmtype, &createtime, &Pfirstname, &Plastname, &Pemail, &Pphone, &Pmobile, &Sfirstname, &Slastname, &Semail, &Sphone, &Smobile)
>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51
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
<<<<<<< HEAD


	if (id!=emp.Device){
		BadResponse(w,"Enter  Correct ID OR Time Range")
	}
	w.Header().Set("Content-Type", "application/json")
//	fmt.Print(res)
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
=======

	w.Header().Set("Content-Type", "application/json")
	fmt.Print(res)
	messagee, err := json.Marshal(res)
>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51
	if err != nil {
		fmt.Println(err)
	}
	w.Write(messagee)
<<<<<<< HEAD
}
=======

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
>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51

func Session() string {
	person := &Person{""}
	buf, _ := json.Marshal(person)
	body := bytes.NewBuffer(buf)
	response, err := http.Post("https://api.obitalk.com/api/v1/sessions", "application/json", body)
	if err != nil {
		fmt.Println("Error")

	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var data TokenData
	json.Unmarshal(responseData, &data)
	token := data.Data.Token
	return token
}
