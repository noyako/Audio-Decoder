package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noyako/Audio-Decoder/storage"
)

type UserService struct {
	db   storage.User
	addr string
}

func (u *UserService) process() {
	mux := mux.NewRouter()
	mux.HandleFunc("/login", u.Login).Methods("POST")
	mux.HandleFunc("/register", u.Register).Methods("POST")
	mux.HandleFunc("/refresh", u.RefreshToken).Methods("POST")

	mux.HandleFunc("/code/all", u.GetAll).Methods("GET")
	mux.HandleFunc("/code/get", u.GetOne).Methods("GET")
	mux.HandleFunc("/code/download", u.Load).Methods("GET")

	mux.HandleFunc("/code/encode", u.Encode).Methods("POST")
	mux.HandleFunc("/code/decode", u.Decode).Methods("POST")

	server := http.Server{
		Handler: mux,
		Addr:    "127.0.0.1:8081",
	}
	log.Fatal(server.ListenAndServe())
}
