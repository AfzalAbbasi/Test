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
	Host         = "45.76.175.38:27017"
	AuthUserName = "admin"
	AuthPassword = "Lmkt@ptcl1234"
	AuthDatabase = "admin"
	Collection   = "CoLLections"
)

type LCRData struct {
	Active bool   `json:"active" bson:"active"`
	Number Number `json:"number" bson:"number"`
}

type LCRDataa struct {
	ID     bson.ObjectId `json:"ID" bson:"_id"`
	Active bool          `json:"active" bson:"active"`
	Number Number        `json:"number" bson:"number"`
}

type Number struct {
	Value int  `json:"value" bson:"value"`
	Rates Rate `json:"rates" bson:"rates"`
}

type Rate struct {
	Interstate    float64 `json:"interstate" bson:"interstate"`
	Intrastate    float64 `json:"intrastate" bson:"intrastate"`
	Indeterminate float64 `json:"indeterminate" bson:"indeterminate"`
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
		Addrs:   []string{"45.76.175.38:27017"},
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
	file, err := os.Open("C:\\Users\\Afzal\\Desktop\\Book1.csv")

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
		interstate, err := strconv.ParseFloat(record[1], 64)
		intrastate, err := strconv.ParseFloat(record[2], 64)
		indeterminate, err := strconv.ParseFloat(record[3], 64)
		err = collection.Insert(&LCRData{true, Number{number, Rate{interstate, intrastate, indeterminate}}})
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
	profile := []LCRDataa{}
	//err := collection.Find(bson.M{}).All(&profile)
	err := collection.Find(bson.M{"number.value": p_id}).Sort("number.rates.intrastate").All(&profile)
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
