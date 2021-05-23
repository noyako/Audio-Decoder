package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/noyako/Audio-Decoder/model"
	"github.com/noyako/Audio-Decoder/service"

	"golang.org/x/crypto/bcrypt"
)

type credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ResponceToken struct {
	Token string `json:"token"`
}

func (u *UserService) Login(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	user, err := u.db.GetByName(creds.Name)
	if err != nil {
		service.ProcessForbidden(w, "Cannot find user")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		service.ProcessForbidden(w, "Wrong password")
		return
	}

	service.ProcessOk(w, ResponceToken{user.Token})
}

func (u *UserService) Register(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		service.ProcessBadFormat(w, "Password in wrong format")
		return
	}

	_, err = u.db.GetByName(creds.Name)
	if err != nil {
		user := model.User{
			Username: creds.Name,
			Password: string(hashedPassword),
			Token:    generateToken(),
		}
		err = u.db.Save(&user)
		if err != nil {
			service.ProcessServerError(w, "Error")
			return
		}

		jsonData, _ := json.Marshal(UserName{user.Username})
		resp, err := http.Post("http://localhost:8082/new", "application/json",
			bytes.NewBuffer(jsonData))

		if err != nil {
			service.ProcessServerError(w, "Error")
			return
		}
		defer resp.Body.Close()
		if resp.Status != "200 OK" {
			service.ProcessServerError(w, "Error")
			return
		}

		service.ProcessOk(w, ResponceToken{user.Token})
	} else {
		service.ProcessServerError(w, "User exists")
	}
}

func (u *UserService) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	user, err := u.db.GetByName(creds.Name)
	if err != nil {
		service.ProcessForbidden(w, "Cannot find user")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		service.ProcessForbidden(w, "Wrong password")
		return
	}

	user.Token = generateToken()
	err = u.db.Save(user)
	if err != nil {
		service.ProcessServerError(w, "Error")
		return
	}

	service.ProcessOk(w, ResponceToken{user.Token})
}
