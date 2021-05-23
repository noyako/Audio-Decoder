package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/noyako/Audio-Decoder/service"
)

func (u *UserService) GetAll(w http.ResponseWriter, r *http.Request) {
	var token Token
	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	user, err := u.db.GetByToken(token.Token)
	if err != nil {
		service.ProcessServerError(w, "User not exists")
	} else {
		jsonData, _ := json.Marshal(UserName{user.Username})
		resp, err := http.Post("http://localhost:8082/all", "application/json",
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

		data, _ := ioutil.ReadAll(resp.Body)
		service.ProcessOkString(w, data)
	}
}
