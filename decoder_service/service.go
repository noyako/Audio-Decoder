package main

import (
	"log"
	"net/http"

	"github.com/noyako/swolf"
)

type DecodeService struct {
	adb  *swolf.Dealer
	addr string
}

func (d *DecodeService) process() {
	mux := http.NewServeMux()
	mux.HandleFunc("/new", d.New)
	mux.HandleFunc("/all", d.GetAll)
	mux.HandleFunc("/get", d.GetOne)
	mux.HandleFunc("/download", d.Load)

	mux.HandleFunc("/encode", d.Encode)
	mux.HandleFunc("/decode", d.Decode)

	err := http.ListenAndServe(d.addr, mux)
	log.Fatal(err)
}
