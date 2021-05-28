package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/noyako/Audio-Decoder/request"
	"github.com/noyako/Audio-Decoder/service"
	"github.com/spf13/viper"
)

const (
	generalError = "Error"
	okResponse   = "200 OK"
)

// GetAll returns all user audios
func (u *UserService) GetAll(w http.ResponseWriter, r *http.Request) {
	var token request.UserToken
	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	user, err := u.db.GetByToken(token.Token)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
	} else {
		jsonData, _ := json.Marshal(request.UserName{Username: user.Username})
		s := decoder(viper.GetString("endpoints.decoder.all"))
		req, err := http.NewRequest("GET", s, bytes.NewBuffer(jsonData))
		req.Header.Set("content-type", "application/json")
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}
		defer resp.Body.Close()
		if resp.Status != okResponse {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}

		data, _ := ioutil.ReadAll(resp.Body)
		service.ProcessOkString(w, data)
	}
}

// GetOne returns information about selected audio
func (u *UserService) GetOne(w http.ResponseWriter, r *http.Request) {
	var token request.UserTokenAudioToken
	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	user, err := u.db.GetByToken(token.Token)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
	} else {
		jsonData, _ := json.Marshal(request.UserNameAudioToken{Username: user.Username, Token: token.Audio})
		req, err := http.NewRequest("GET", decoder(viper.GetString("endpoints.decoder.one")), bytes.NewBuffer(jsonData))
		req.Header.Set("content-type", "application/json")
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}
		defer resp.Body.Close()
		if resp.Status != okResponse {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}

		data, _ := ioutil.ReadAll(resp.Body)
		service.ProcessOkString(w, data)
	}
}

// Load dowloads audio
func (u *UserService) Load(w http.ResponseWriter, r *http.Request) {
	var token request.UserTokenAudioToken
	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	user, err := u.db.GetByToken(token.Token)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
	} else {
		jsonData, _ := json.Marshal(request.UserNameAudioToken{Username: user.Username, Token: token.Audio})
		req, err := http.NewRequest("GET", decoder(viper.GetString("endpoints.decoder.load")), bytes.NewBuffer(jsonData))
		req.Header.Set("content-type", "application/json")
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}
		defer resp.Body.Close()
		if resp.Status != okResponse {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}

		io.Copy(w, resp.Body)
		service.ProcessOkFile(w)
	}
}

// Encode encrypts input audio
func (u *UserService) Encode(w http.ResponseWriter, r *http.Request) {
	var token request.UserTokenAudio
	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	if token.AudioInit.URL == "" {
		service.ProcessBadFormat(w, service.ErrAudioURLEmpty)
		return
	}
	if token.AudioInit.Key == "" {
		service.ProcessBadFormat(w, service.ErrEncryptionKeyEmpty)
		return
	}
	if token.AudioInit.KID == "" {
		service.ProcessBadFormat(w, service.ErrEncryptionKeyIDEmpty)
		return
	}

	user, err := u.db.GetByToken(token.Token)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
	} else {
		jsonData, _ := json.Marshal(request.UserAudio{Username: user.Username, AudioInit: token.AudioInit})
		resp, err := http.Post(decoder(viper.GetString("endpoints.decoder.encrypt")), "application/json",
			bytes.NewBuffer(jsonData))

		if err != nil {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}
		defer resp.Body.Close()
		if resp.Status != okResponse {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}

		data, _ := ioutil.ReadAll(resp.Body)
		service.ProcessOkString(w, data)
	}
}

// Decode decrypts audio
func (u *UserService) Decode(w http.ResponseWriter, r *http.Request) {
	var token request.UserTokenAudio
	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	if token.AudioInit.URL == "" {
		service.ProcessBadFormat(w, service.ErrAudioURLEmpty)
		return
	}
	if token.AudioInit.Key == "" {
		service.ProcessBadFormat(w, service.ErrEncryptionKeyEmpty)
		return
	}

	user, err := u.db.GetByToken(token.Token)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
	} else {
		jsonData, _ := json.Marshal(request.UserAudio{Username: user.Username, AudioInit: token.AudioInit})
		resp, err := http.Post(decoder(viper.GetString("endpoints.decoder.decrypt")), "application/json",
			bytes.NewBuffer(jsonData))

		if err != nil {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}
		defer resp.Body.Close()
		if resp.Status != okResponse {
			log.Println(err)
			service.ProcessServerError(w, generalError)
			return
		}

		data, _ := ioutil.ReadAll(resp.Body)
		service.ProcessOkString(w, data)
	}
}

func decoder(endpoint string) string {
	return "http://" + viper.GetString("decoder.service.url") + ":" + viper.GetString("decoder.service.port") + endpoint
}
