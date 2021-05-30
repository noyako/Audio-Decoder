package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noyako/swolf"
	"github.com/spf13/viper"
)

// DecodeService represents service to work with audio
type DecodeService struct {
	adb  *swolf.Dealer
	addr string
}

func (d *DecodeService) process() {
	mux := mux.NewRouter()

	mux.HandleFunc(viper.GetString("endpoints.decoder.new"), d.New).Methods("POST")
	mux.HandleFunc(viper.GetString("endpoints.decoder.all"), d.GetAll).Methods("GET")
	mux.HandleFunc(viper.GetString("endpoints.decoder.one"), d.GetOne).Methods("GET")
	mux.HandleFunc(viper.GetString("endpoints.decoder.load"), d.Load).Methods("GET")
	mux.HandleFunc(viper.GetString("endpoints.decoder.delete"), d.Remove).Methods("POST")

	mux.HandleFunc(viper.GetString("endpoints.decoder.encrypt"), d.Encode).Methods("POST")
	mux.HandleFunc(viper.GetString("endpoints.decoder.decrypt"), d.Decode).Methods("POST")

	server := http.Server{
		Handler: mux,
		Addr:    d.addr,
	}
	log.Fatal(server.ListenAndServe())
}
