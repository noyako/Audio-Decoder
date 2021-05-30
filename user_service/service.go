package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noyako/Audio-Decoder/storage"
	"github.com/spf13/viper"
)

// UserService represents a service to work with users
type UserService struct {
	db   storage.User
	addr string
}

func (u *UserService) process() {
	mux := mux.NewRouter()
	mux.HandleFunc(viper.GetString("endpoints.user.login"), u.Login).Methods("POST")
	mux.HandleFunc(viper.GetString("endpoints.user.register"), u.Register).Methods("POST")
	mux.HandleFunc(viper.GetString("endpoints.user.refresh"), u.RefreshToken).Methods("POST")

	mux.HandleFunc(viper.GetString("endpoints.user.all"), u.GetAll).Methods("GET")
	mux.HandleFunc(viper.GetString("endpoints.user.one"), u.GetOne).Methods("GET")
	mux.HandleFunc(viper.GetString("endpoints.user.load"), u.Load).Methods("GET")
	mux.HandleFunc(viper.GetString("endpoints.user.delete"), u.Remove).Methods("POST")

	mux.HandleFunc(viper.GetString("endpoints.user.encrypt"), u.Encode).Methods("POST")
	mux.HandleFunc(viper.GetString("endpoints.user.decrypt"), u.Decode).Methods("POST")

	server := http.Server{
		Handler: mux,
		Addr:    u.addr,
	}
	log.Fatal(server.ListenAndServe())
}
