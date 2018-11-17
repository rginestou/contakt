package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Config ...
type Config struct {
	Port     int
	DBServer string
	DBName   string
	Key      string
	Login    string
	Password string
}

var config Config
var store *sessions.CookieStore

func router() *mux.Router {
	r := mux.NewRouter()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	fs = http.FileServer(http.Dir("data"))
	r.PathPrefix("/data/").Handler(http.StripPrefix("/data/", fs))

	// Auth
	r.HandleFunc("/login", loginGET).Methods("GET")
	r.HandleFunc("/login", loginPOST).Methods("POST")
	r.HandleFunc("/logout", logout)

	// Routes
	r.HandleFunc("/contact", sessionAuth(contactGET)).Methods("GET")
	r.HandleFunc("/contact", sessionAuth(contactPOST)).Methods("POST")
	r.HandleFunc("/edit", sessionAuth(editContactGET)).Methods("GET")
	r.HandleFunc("/new", sessionAuth(newContactGET)).Methods("GET")
	r.HandleFunc("/delete", sessionAuth(deleteContactPOST)).Methods("POST")
	r.HandleFunc("/", sessionAuth(indexGET)).Methods("GET")

	return r
}

func main() {
	// Parse config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}

	// Auth
	store = sessions.NewCookieStore([]byte(config.Key))

	// Database
	connect(config.DBServer, config.DBName)

	// Serve
	port := strconv.Itoa(config.Port)
	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router()))
}
