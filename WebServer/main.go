package main

import (
	"encoding/json"
	"net/http"
)

type name struct {
	First string `json:"Name"`
}

func main() {
	http.HandleFunc("/user", json_responce)
	http.ListenAndServe(":8080", nil)

}
func json_responce(w http.ResponseWriter, r *http.Request) {
	person := name{"Afzal"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)

}
