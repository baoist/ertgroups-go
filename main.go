package main

import (
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"github.com/baoist/ertgroups/encoder"
)

type CsvRequest struct {
	Body string `json:"csv"`
	ColDelimRe string `json:"col_delimiter_regex"`
	ColExtraRe string `json:"col_extra_character_regex"`
	RowDelimRe string `json:"row_delimiter"`
}

type JsonErr struct {
	Message string `json:"message"`
}


func JSONError(w http.ResponseWriter, err string, code int) {
	json_error := JsonErr{
		Message: err,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(json_error)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("pong"))
}

func groupsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		JSONError(w, "Invalid method in request. Expects `POST`.", 405)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	csv := CsvRequest{
		Body: "",
		ColDelimRe: `,`,
		ColExtraRe: `\W.*`,
		RowDelimRe: `\n`,
	}
	json.Unmarshal(reqBody, &csv)

	delimiters := encoder.Delimiters{
		ColDelimRe: csv.ColDelimRe,
		ColExtraRe: csv.ColExtraRe,
		RowDelimRe: csv.RowDelimRe,
	}
	group := encoder.CreateGroup(csv.Body, delimiters)
	group.Format()

	js, err := json.Marshal(group)
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

	w.Write([]byte(js))
}

func handleRequests() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/create", groupsCreate)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
