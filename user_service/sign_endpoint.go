package main

import (
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		service.ProcessBadFormat(w, "Password in wrong format")
		return
	}

	user, err := u.db.GetByCreds(creds.Name, string(hashedPassword))
	if err != nil {
		service.ProcessForbidden(w, "Cannot find user")
		return
	}

	service.ProcessOk(w, user)
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

	_, err = u.db.GetByCreds(creds.Name, string(hashedPassword))
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		service.ProcessBadFormat(w, "Password in wrong format")
		return
	}

	user, err := u.db.GetByCreds(creds.Name, string(hashedPassword))
	if err != nil {
		service.ProcessForbidden(w, "Cannot find user")
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
