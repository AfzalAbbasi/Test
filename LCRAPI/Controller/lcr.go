package Controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AfzalAbbasi/Test/LCRAPI/Model"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		session *mgo.Session
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(a *mgo.Session) *UserController {
	return &UserController{a}
}

func (uc UserController) GetUser(c echo.Context) error {
	//session.Copy()
	collection := uc.session.DB("lcr").C("test")
	p_id, _ := strconv.Atoi(c.Param("id"))
	profile := []Model.LCRDataa{}
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
func (uc UserController) Carrier(c echo.Context) error {

	collection := uc.session.DB("lcr").C("carrier")
	// Destination

	err := collection.Insert(&Model.Carrier{"tets", 50, true, time.Now().UTC(), time.Now().UTC()})
	if err != nil {
		CreateBadResponse(&c, http.StatusNotFound, "Data", "Not Found")
	}

	return c.JSON(http.StatusOK, "Data Are uploaded")
}
func (uc UserController) Postcsv(c echo.Context) error {
	collection := uc.session.DB("lcr").C("test")
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
		record := data
		if err == io.EOF {
			//os.Remove(path)
			return c.JSON(http.StatusOK, "Data are Uploaded")
		}
		number, err := strconv.Atoi(record[0])
		interstate, err := strconv.ParseFloat(record[1], 64)
		intrastate, err := strconv.ParseFloat(record[2], 64)
		indeterminate, err := strconv.ParseFloat(record[3], 64)
		err = collection.Insert(&Model.LCRData{p_id, true, Model.Number{number, Model.Rate{interstate, intrastate, indeterminate}}})
		if err != nil {
			fmt.Println(err)
			CreateBadResponse(&c, http.StatusNotFound, "Data", "Not Found")
		}

	}
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
