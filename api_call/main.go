package main

import (
"fmt"
"net/http"
"regexp"

"github.com/labstack/echo"
"github.com/labstack/echo/middleware"
"github.com/afzalabbasi/testrepo/fs"
)

type Message struct {
	Message string `json:"message"`
}
type Call struct {
	To_DID        string `json:"to_did"`
	From_DID      string `json:"from_did"`
	Tx_DID        string `json:"tx_did"`
	Human_audio   string `json:"human_audio"`
	Vm_detect     bool   `json:"vm_detect"`
	Vm_drop       bool   `json:"vm_drop"`
	Vm_audio      string `json:"vm_audio"`
	Compaign_type int    `json:"compaign_type"`
	Tx_DTMF       int    `json:"tx_dtmf"`
	DND_DTMF      int    `json:"dnd_dtmf"`
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.POST("/call", callapi)

	e.Logger.Fatal(e.Start(":8080"))

}
func callapi(c echo.Context) error {
	u := new(Call)
	if err := c.Bind(u); err != nil {
		return err
	}
	valid_Tx := IsUSNumber(u.Tx_DID)
	valid_To := IsUSNumber(u.To_DID)
	valid_from := IsUSNumber(u.From_DID)
	r, err1 := regexp.Compile(`^http://`)
	if err1 != nil {
		fmt.Println(err1)
	}

	r1, err1 := regexp.Compile("\\.mp3$")
	if err1 != nil {
		fmt.Println(err1)
	}
	valid_Haudio := r1.MatchString(u.Human_audio)
	valid_Vaudio := r1.MatchString(u.Vm_audio)
	valid_Haudioo := r.MatchString(u.Human_audio)
	valid_Vaudioo := r.MatchString(u.Vm_audio)
	if valid_Haudioo == false {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct Human_Audio URL"})
	}
	if valid_Vaudioo == false {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct Vm_Audio URL"})
	}
	if valid_Haudio == false {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct Human_Audio URL"})
	}
	if valid_Vaudio == false {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct Vm_Audio URL"})
	}
	if valid_Tx == false {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct Tx_DID"})
	}
	if valid_To == false {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct To_DID"})
	}
	if valid_from == false {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct From_DID"})
	}
	if u.DND_DTMF < 0 || u.DND_DTMF >= 10 {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct DND_DMTF"})
	}
	if u.Tx_DTMF < 0 || u.Tx_DTMF >= 10 {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct TX_DTMF"})
	}
	if u.Compaign_type < 0 || u.Compaign_type >= 10 {
		return c.JSON(http.StatusNotFound, Message{"Please Enter Correct Compaign Type"})
	}
	if(u.Tx_DTMF==u.DND_DTMF){
		return c.JSON(http.StatusNotFound, Message{"TX_DTMF and DND_DTMF could not be the same"})
	}
	fs.CallStart(Call{u.Tx_DID,u.From_DID,u.Tx_DID,u.Human_audio,u.Vm_detect,u.Vm_drop,u.Vm_audio,u.Compaign_type,u.DND_DTMF,u.Tx_DTMF})
		return c.JSON(http.StatusCreated, u)
}
func IsUSNumber(number string) bool {
	match, err := regexp.MatchString("^(?:(?:\\+?1\\s*(?:[.-]\\s*)?)?(?:\\(\\s*([2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9])\\s*\\)|([2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9]))\\s*(?:[.-]\\s*)?)([2-9]1[02-9]|[2-9][02-9]1|[2-9][02-9]{2})\\s*(?:[.-]\\s*)?([0-9]{4})(?:\\s*(?:#|x\\.?|ext\\.?|extension)\\s*(\\d+))?$", number)
	if err != nil {
		return false
	}
	return match
}
