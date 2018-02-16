package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type name struct {
	Name string `json:"Message "`
}

func main() {
	http.HandleFunc("/", URLValues)
	http.HandleFunc("/user", From_File_Upload)
	http.ListenAndServe(":8080", nil)
}

func URLValues(w http.ResponseWriter, req *http.Request) {

	v := req.FormValue("message")
	w.Header().Set("Content-Type", "application/json")
	p2 := name{v}

	if v == "bilal" {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	json.NewEncoder(w).Encode(p2)

}
func From_File_Upload(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		key := "my_file"
		file, _, err := req.FormFile(key)
		if err != nil {

			http.Error(w, err.Error(), 500)
			return
		}
		defer file.Close()
		src := io.LimitReader(file, 400)
		dst, err := os.Create(filepath.Join(".", "file.txt"))
		if err != nil {

			http.Error(w, err.Error(), 500)
			return
		}
		defer dst.Close()
		io.Copy(dst, src)
	}

	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, `<form method="POST" enctype="multipart/form-data">
		<input type="file" name="my_file">
         <input type="submit">
          </form>`)
}
