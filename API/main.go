package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonlvhit/gocron"
	"golang.org/x/time/rate"
	"gopkg.in/mgo.v2"
)

var limiter = rate.NewLimiter(2, 5)

type State struct {
	DisplayString string `json:"displayString"`
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
	DeviceData Device `json:"data"`
}
type SessionData struct {
	Token     string `json:"token"`
	Agent     string `json:"agent"`
	Org       string `json:"org"`
	CreateAt  string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
}

type Device []struct {
	ID           string `json:"id" bson:"id"`
	ObiNumber    string `json:"obiNumber" bson:"obiNumber"`
	MacAddress   string `json:"macAddress" bson:"macAddress"`
	SerialNumber string `json:"serialNumber" bson:"serialNumber"`
	Org          string `json:"org"  bson:"org"`
	DeviceType   string `json:"deviceType" bson:"deviceType"`
}

var db *mgo.Database
var Collections *mgo.Collection
var (
	session    *mgo.Session
	collection *mgo.Collection
)

//const MongoDb details
const (
	Host         = ""
	AuthUserName = "admin"
	AuthPassword = ""
	AuthDatabase = "admin"
	Collection   = "CoLLections"
)

func main() {
	
	Devices()

}
func Devices() {
	//	session.Copy()
	//defer session.Close()
	token := Session()
	fmt.Println(token)
	value := fmt.Sprintf("%s %s", "Bearer", token)

	client := &http.Client{}
	request, _ := http.NewRequest("GET", "", nil)
	request.Header.Set("authorization", value)
	response, err := client.Do(request)

	if err != nil {
		fmt.Print(err.Error())
	}
	req1, err := ioutil.ReadAll(response.Body)
	var devicedata DeviceData
	json.Unmarshal(req1, &devicedata)
	for _, item := range devicedata.DeviceData {
		a := item.ID

		api := fmt.Sprintf("%s/%s/%s", "", a, "quick-values")
		request1, _ := http.NewRequest("POST", api, nil)
		request1.Header.Set("authorization", value)
		response1, err := client.Do(request1)

		if err != nil {
			fmt.Print(err.Error())
		}
		var data StatusData
		req, err := ioutil.ReadAll(response1.Body)
		json.Unmarshal(req, &data)

		var b []State
		for _, item := range data.Data.SpRegStates {
			a := item.DisplayString
			if a != "" {
				if a != "Service Not Configured" && a != "null" {
					b = append(b, State{DisplayString: a})

				}
			}
		}
		db, err := sql.Open("mysql", "")
		if err != nil {
			fmt.Println(err)
		}
		data.Data.SpRegStates = b
		i := fmt.Sprintf("%s", b)

		_, err = db.Query("INSERT INTO Devices(device,online,firmware,spRegState) VALUES (?,?,?,?)", data.Data.Device, data.Data.Online, data.Data.Firmware, i)
		//insert, err := db.Query("INSERT  into Test VALUES (?,?,?,?,?,?,?)")
		if err != nil {
			fmt.Println(err)
		}
		//	fmt.Println("Uploaded")
		defer db.Close()
		//collection = session.DB("lcr").C("db")
		//err = collection.Insert(data)

	}
	gocron.Every(1).Minute().Do(Devices)
	<-gocron.Start()
}
func Session() string {
	person := &Person{"email", "pass}
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
