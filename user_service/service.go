package main

import (
	"log"
	"net/http"

	"github.com/noyako/Audio-Decoder/storage"
)

type UserService struct {
	db   storage.User
	addr string
}

func (u *UserService) process() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", u.Login)
	mux.HandleFunc("/register", u.Register)
	mux.HandleFunc("/refresh", u.RefreshToken)

	mux.HandleFunc("/code/all", u.GetAll)

	err := http.ListenAndServe(u.addr, mux)
	log.Fatal(err)
}
