package main

import (
	"fmt"
	"regexp"
<<<<<<< HEAD
	//	"time"
	"github.com/zemirco/couchdb"
	"log"
	"net/url"
	"strings"
	"time"
=======
	"github.com/zemirco/couchdb"
	//"github.com/fjl/go-couchdb"

	"net/url"

	"log"


>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51
)

type Data struct {
	Name string `json:"name"`
}

const (
<<<<<<< HEAD
	Host = ""
)

func main() {

	fmt.Print("\n")
	_ = time.Now().UTC()
	time.Now().AddDate(0, 0, 0)
	PrevMonth := time.Now().UTC().Add(-12960 * time.Hour)
	//fmt.Println(PrevMonth)
	fmt.Print("\n")
=======
	Host = ""
)

func main() {
>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51
	u, err := url.Parse(Host)
	if err != nil {
		fmt.Print(err)
	}
	client, err := couchdb.NewClient(u)
	if err != nil {
		log.Print(err)
	}
	var data []string
<<<<<<< HEAD
	data, err = client.All()
	if err != nil {
		log.Print(err)
	}
	for _, item := range data {

		matched, err := regexp.MatchString("-", item)
		if err != nil {
			fmt.Println(err)
		}
		if matched == true {
			s := strings.Split(item, "-")
			datee := s[1]
			month := fmt.Sprint(datee, "01")
			const longForm = "20060102"
			db_Date, _ := time.Parse(longForm, month)
			if db_Date.After(PrevMonth) {

			} else {
				client.Delete(item)

			}

		}

	}
=======
	data, err=client.All()
	if err != nil {
		log.Print(err)
	}
	for _, item :=range data{
		matched, err := regexp.MatchString("-", item)
		if err!=nil{
			fmt.Println(err)
		}
		if matched==true{
			fmt.Println(item)
		}

	}


	}


>>>>>>> ce325fd8a343c52aaca454d25d65dbbd2380ea51

