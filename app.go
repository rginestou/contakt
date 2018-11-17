package main

import (
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
)

// Config ...
type Config struct {
	Port     int
	Secret   string
	DBServer string
	DBName   string
}

var config Config

func router() *mux.Router {
	r := mux.NewRouter()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Routes
	r.HandleFunc("/contact", contactHandler)
	r.HandleFunc("/edit", editContactGET).Methods("GET")
	r.HandleFunc("/new", newContactGET).Methods("GET")
	r.HandleFunc("/delete", deleteContactPOST).Methods("POST")
	r.HandleFunc("/", indexGET).Methods("GET")

	return r
}

func main() {
	// Parse config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}

	connect(config.DBServer, config.DBName)

	// Serve
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router()))
}
