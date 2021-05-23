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
	// mux.HandleFunc("/status", u.Status)

	err := http.ListenAndServe(d.addr, mux)
	log.Fatal(err)
}
