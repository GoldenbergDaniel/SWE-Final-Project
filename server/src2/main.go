package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "data.db")
	db.Exec("")

	r := mux.NewRouter()
	r.HandleFunc("/", PostLogin).Methods("POST")

	fmt.Println("Listening on localhost:5174")
	http.ListenAndServe(":5174", r)
}

func PostLogin(res http.ResponseWriter, req *http.Request) {
	SetCors(res.Header())

	vars := mux.Vars(req)
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "Category: %v\n", vars["category"])
}

func SetCors(header http.Header) {
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH")
	header.Set("Access-Control-Allow-Headers", "Content-Type")
}
