package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var path = "F:\\file.csv"
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

type LCRData struct {
	CarrierID bson.ObjectId `json:"carrierid" bson:"_carrierid"`
	Active    bool          `json:"active" bson:"active"`
	Number    Number        `json:"number" bson:"number"`
}
type Carrier struct {
	//ID          bson.ObjectId `json:"ID" bson:"_id"`
	Name        string
	DispatcerID int
	Active      bool
	CreateDate  time.Time
	UpdateDate  time.Time
}

type LCRDataa struct {
	//ID          bson.ObjectId `json:"ID" bson:"_id"`
	ID        bson.ObjectId `json:"ID" bson:"_id"`
	CarrierID bson.ObjectId `json:"carrierid " bson:"_carrierid"`
	Active    bool          `json:"active" bson:"active"`
	Number    Number        `json:"number" bson:"number"`
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
	e.Use(middleware.Static("Template"))
	e.POST("/users/:id", postcsv)
	e.POST("/carrier", carrier)
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

	session.Copy()
	collection = session.DB("lcr").C("test")
	p_id := bson.ObjectIdHex(c.Param("id"))
	// Destination
	key := "myfile"
	file, err := c.FormFile(key)
	if err != nil {
		//	return err
		CreateBadResponse(&c, http.StatusNotFound, "Data", "Not Found")
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusOK, "Not found")
		return err
	}
	defer src.Close()
	dst, err := os.Create("F:\\file.csv")
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	defer dst.Close()
	Csvfile, err := os.Open("F:\\file.csv")
	defer Csvfile.Close()
	reader := csv.NewReader(Csvfile)

	for {
		data, err := reader.Read()
		record = data
		if err == io.EOF {
			//os.Remove(path)
			return c.JSON(http.StatusOK, "Data are Uploaded")
		}
		number, err := strconv.Atoi(record[0])
		interstate, err := strconv.ParseFloat(record[1], 64)
		intrastate, err := strconv.ParseFloat(record[2], 64)
		indeterminate, err := strconv.ParseFloat(record[3], 64)
		err = collection.Insert(&LCRData{p_id, true, Number{number, Rate{interstate, intrastate, indeterminate}}})
		if err != nil {
			fmt.Println(err)
			CreateBadResponse(&c, http.StatusNotFound, "Data", "Not Found")
		}

	}
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
	return CreateSuccessResponse(&c, 200, "Get Data", "Successful", b)

}
func carrier(c echo.Context) error {
	session.Copy()
	collection = session.DB("lcr").C("carrier")
	// Destination

	err := collection.Insert(&Carrier{"tets", 50, true, time.Now().UTC(), time.Now().UTC()})
	if err != nil {
		CreateBadResponse(&c, http.StatusNotFound, "Data", "Not Found")
	}

	return c.JSON(http.StatusOK, "Data Are uploaded")
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
