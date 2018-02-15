package main

import (
	"html/template"
	"log"
	"net/http"
)

type person struct {
	Name string
	Last string
}

func main() {

	tpl, err := template.ParseFiles("form.html")
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fName := r.FormValue("first")
		lname := r.FormValue("last")
		err = tpl.Execute(w, person{fName, lname})
		if err != nil {

			http.Error(w, err.Error(), 500)
			log.Fatalln(err)
		}
	})

	http.ListenAndServe(":8080", nil)

}
