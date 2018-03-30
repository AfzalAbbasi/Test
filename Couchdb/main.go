package main

import (
	"fmt"
	"regexp"
	"github.com/zemirco/couchdb"
	//"github.com/fjl/go-couchdb"

	"net/url"

	"log"


)

type Data struct {
	Name string `json:"name"`
}

const (
	Host = "http://dev.venturetel.co:15984"
)

func main() {
	u, err := url.Parse(Host)
	if err != nil {
		fmt.Print(err)
	}
	client, err := couchdb.NewClient(u)
	if err != nil {
		log.Print(err)
	}
	var data []string
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



