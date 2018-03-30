package main

import (
	// Standard library packages
	"log"
	"net/http"
	"time"

	// Third party packages
	"github.com/AfzalAbbasi/Test/LCRAPI/Controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
)

func main() {
	// Instantiate a new router
	uc := Controller.NewUserController(getSession())
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	//e.Use(middleware.Static("Template"))
	e.POST("/users/:id", uc.Postcsv)
	e.POST("/carrier", uc.Carrier)
	e.GET("/users/:id", uc.GetUser)
	// Get a UserController instance

	// Get a user resource

	// Fire up the server
	http.ListenAndServe("localhost:8080", e)
}
func getSession() *mgo.Session {
	// Connect to our local mongo
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:   []string{""},
		Timeout: 60 * time.Second,
		//Database: AuthDatabase,
		//Username: AuthUserName,
		//Password: AuthPassword,
	}

	a, err := mgo.DialWithInfo(mongoDBDialInfo)

	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	// Optional. Switch the session to a monotonic behavior.
	a.SetMode(mgo.Monotonic, true)

	// Deliver session
	return a

}
