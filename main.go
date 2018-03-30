package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
type Message struct {
	message string `json:"message"`
}
type Delete struct {
	delete string `json:"delete"`
}
type Data struct {
	data string `json:"data"`
	User
}

var (
	counter = 0
	name    string
	last    string
	e       = echo.New()
)

func main() {
	// Echo instance

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Static("/", "public")
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/update/:id", updateUser)
	e.DELETE("/delete/:id", deleteUser)
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler

func createUser(c echo.Context) error {

	// Read form fields
	name = c.FormValue("first")
	last = c.FormValue("last")
	counter++
	data := &User{counter, name, last}

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	if err != nil {

		return err
	}
	fileName := fmt.Sprintf("./%d.json", counter)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {

		return err
	}
	f.Write(b)
	f.Close()
	return CreateSuccessResponse(&c, 200, "successful", "data", b)
	//return c.JSON(http.StatusOK, User{counter, name, last})
}

func getUser(c echo.Context) error {
	count, _ := strconv.Atoi(c.Param("id"))
	fileName := fmt.Sprintf("./%d.json", count)
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var data User
	json.Unmarshal(raw, &data)
	fmt.Println(data)

	return c.JSON(http.StatusOK, &data)

}

func deleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	fileName := fmt.Sprintf("./%d.json", id)
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		fmt.Print("error")
	}
	var data User
	json.Unmarshal(raw, &data)
	os.Remove(fileName)

	return c.JSON(http.StatusOK, &Delete{"Deleted"})
}
func updateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	fileName := fmt.Sprintf("./%d.json", id)
	var data = User{counter, "danish", "khan"}
	b, err := json.Marshal(data)
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0666)
	if err != nil {

		return err
	}
	f.Write(b)
	f.Close()
	return c.JSON(http.StatusOK, &Message{"User are Updated"})
}
func CreateSuccessResponse(c *echo.Context, requestCode int, message string, subMessage string, data []byte) error {

	localC := *c
	response := fmt.Sprintf("{\"data\":%s,\"message\":%q,\"submessage\":%q}", data, message, subMessage)
	fmt.Print(response)
	return localC.JSONBlob(requestCode, []byte(response))
}
