package main

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
	"log"
	"net/http"
)

var (
	Db    *memdb.MemDB
	Users []User
	cart  Cart
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/checkout", checkout).Methods("POST")
	myRouter.HandleFunc("/add/{productId}", addToCart).Methods("POST")
	myRouter.HandleFunc("/login", login).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/logout", logout).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	Db = setupDb()
	handleRequests()
}