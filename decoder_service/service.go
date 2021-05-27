package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noyako/swolf"
)

type DecodeService struct {
	adb  *swolf.Dealer
	addr string
}

func (d *DecodeService) process() {
	mux := mux.NewRouter()

	mux.HandleFunc("/new", d.New).Methods("POST")
	mux.HandleFunc("/all", d.GetAll).Methods("GET")
	mux.HandleFunc("/get", d.GetOne).Methods("GET")
	mux.HandleFunc("/download", d.Load).Methods("GET")

	mux.HandleFunc("/encode", d.Encode).Methods("POST")
	mux.HandleFunc("/decode", d.Decode).Methods("POST")

	server := http.Server{
		Handler: mux,
		Addr:    "127.0.0.1:8082",
	}
	log.Fatal(server.ListenAndServe())
}
