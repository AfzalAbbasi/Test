package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Contracttype struct {
	Id    bson.ObjectId `json:"id"   bson:"_id,omitempty"`
	Type1 int           `json:"type" bson:"type" `
	Name  string        `json:"name" bson:"name"  `
	Desc  string        `json:"desc"  bson:"desc"  `
}

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
	e.GET("/users/:id", getUser)
	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{""},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	collection = session.DB("callmylistdb").C("contracttype")
	e.Logger.Fatal(e.Start(":8080"))

}
func getUser(c echo.Context) error {
	count, _ := strconv.Atoi(c.Param("id"))
	profile := []Contracttype{}
	//err := collection.Find(bson.M{}).All(&profile)
	err := collection.Find(bson.M{"type": count}).All(&profile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return c.JSON(http.StatusOK, profile)

}
