package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	session    *mgo.Session
	collection *mgo.Collection
)
var record []string

//const MongoDb details
const (
	Host         = ""
	AuthUserName = "admin"
	AuthPassword = ""
	AuthDatabase = "admin"
	Collection   = "CoLLections"
)

type Mongo struct {
	PreferredPhone int    `json:"preferred_phone" bson:"preferred_phone"`
	PlayerID       int    `json:"player_id" bson:"player_id"`
	FirstName      string `json:"first_name" bson:"fname"`
	LastName       string `json:"last_name" bson:"lname"`
	BirthDay       string `json:"birth_day" bson:"bday"`
}

var db *mgo.Database
var Collections *mgo.Collection

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	//e.Static("/", "public")
	e.POST("/users", postcsv)
	e.GET("/users/:id", getUser)
	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:   []string{""},
		Timeout: 60 * time.Second,
		//Database: AuthDatabase,
		//Username: AuthUserName,
		//Password: AuthPassword,
	}

	a, err := mgo.DialWithInfo(mongoDBDialInfo)
	session = a
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	e.Logger.Fatal(e.Start(":8080"))

}
func postcsv(c echo.Context) error {
	//	var b []byte
	session.Copy()
	collection = session.DB("lcr").C("test")
	file, err := os.Open("C:\\Users\\Afzal\\Desktop\\Test.csv")

	if err != nil {
		CreateBadResponse(&c, http.StatusNotFound, "Data", "Not Found")
	}
	defer file.Close()
	reader := csv.NewReader(file)

	for {
		data, err := reader.Read()
		record = data
		//a :=[]Mongo{}

		if err == io.EOF {
			//return CreateSuccessResponse(&c, http.StatusOK, "Data", "are Uploaded", record)
			return c.JSON(http.StatusOK, "Data are Uploaded")
		} else if err != nil {
			return c.JSON(http.StatusNotFound, err)
		}
		number, err := strconv.Atoi(record[0])
		ID, err := strconv.Atoi(record[1])

		err = collection.Insert(&Mongo{number, ID, record[2], record[3], record[4]})

		if err != nil {
			CreateBadResponse(&c, http.StatusNotFound, "Data", "Not Found")
		}

	}

	//defer session.Close()
	return c.JSON(http.StatusOK, record)

}
func getUser(c echo.Context) error {
	session.Copy()
	collection = session.DB("lcr").C("test")
	p_id, _ := strconv.Atoi(c.Param("id"))
	profile := []Mongo{}
	//err := collection.Find(bson.M{}).All(&profile)
	err := collection.Find(bson.M{"player_id": p_id}).All(&profile)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	b, err := json.Marshal(profile)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	//defer session.Close()
	return CreateSuccessResponse(&c, 200, "Get Data", "Successful", b)

}

func CreateSuccessResponse(c *echo.Context, requestCode int, message string, subMessage string, data []byte) error {
	localC := *c
	response := fmt.Sprintf("{\"data\":%s,\"message\":%q,\"submessage\":%q}", data, message, subMessage)
	fmt.Print(response)
	return localC.JSONBlob(requestCode, []byte(response))
}
func CreateBadResponse(c *echo.Context, requestCode int, message string, subMessage string) error {
	localC := *c
	response := fmt.Sprintf("{\"data\":{},\"message\":%q,\"submessage\":%q}", message, subMessage)
	return localC.JSONBlob(requestCode, []byte(response))
}
