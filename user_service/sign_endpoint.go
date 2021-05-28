package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/noyako/Audio-Decoder/model"
	"github.com/noyako/Audio-Decoder/request"
	"github.com/noyako/Audio-Decoder/service"
	"github.com/spf13/viper"

	"golang.org/x/crypto/bcrypt"
)

type credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Login returns token by credentials
func (u *UserService) Login(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	user, err := u.db.GetByName(creds.Name)
	if err != nil {
		service.ProcessForbidden(w, service.ErrFindUser)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		service.ProcessForbidden(w, service.ErrFindUser)
		return
	}

	service.ProcessOk(w, request.UserToken{Token: user.Token})
}

// Register new user (tenant)
func (u *UserService) Register(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	_, err = u.db.GetByName(creds.Name)
	if err != nil {
		user := model.User{
			Username: creds.Name,
			Password: string(hashedPassword),
			Token:    service.GenerateToken(),
		}
		err = u.db.Save(&user)
		if err != nil {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}

		jsonData, _ := json.Marshal(request.UserName{Username: user.Username})
		resp, err := http.Post(decoder(viper.GetString("endpoints.decoder.new")), "application/json",
			bytes.NewBuffer(jsonData))

		if err != nil {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}
		defer resp.Body.Close()
		if resp.Status != okResponse {
			log.Println(resp)
			service.ProcessServerError(w, generalError)
			return
		}

		service.ProcessOk(w, request.UserToken{Token: user.Token})
	} else {
		service.ProcessServerError(w, service.ErrFindUser)
	}
}

// RefreshToken generates new token for user
func (u *UserService) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	user, err := u.db.GetByName(creds.Name)
	if err != nil {
		service.ProcessForbidden(w, service.ErrFindUser)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		service.ProcessForbidden(w, service.ErrFindUser)
		return
	}

	user.Token = service.GenerateToken()
	err = u.db.Save(user)
	if err != nil {
		log.Println(err)
		service.ProcessServerError(w, generalError)
		return
	}

	service.ProcessOk(w, request.UserToken{Token: user.Token})
}
