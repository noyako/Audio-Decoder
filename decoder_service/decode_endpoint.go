package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/noyako/Audio-Decoder/model"
	"github.com/noyako/Audio-Decoder/request"
	"github.com/noyako/Audio-Decoder/service"
	"github.com/noyako/Audio-Decoder/storage"
)

const (
	statusDone    = "done"
	statusProcess = "processing"
	statusError   = "error"

	encryptExtension = ".crypt"
	audioExtension   = ".mp3"
)

// New creates new tenant
func (d *DecodeService) New(w http.ResponseWriter, r *http.Request) {
	var req request.UserName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	db, err := d.adb.Create(req.Username)
	if err != nil {
		service.ProcessServerError(w, service.ErrCreateUser)
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, service.ErrCreateUser)
		return
	}

	as.Migrate()

	os.MkdirAll(getBaseDir(req.Username), 0777)
}

// GetAll returns all published audios for tenant
func (d *DecodeService) GetAll(w http.ResponseWriter, r *http.Request) {
	var req request.UserName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}
	as.Migrate()

	audios, err := as.GetAll()
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUserAudio)
		return
	}

	var statuses []request.AudioStatus
	for i := range audios {
		statuses = append(statuses, request.AudioStatus{
			Name:  audios[i].Name,
			Date:  audios[i].PostedAt,
			Token: audios[i].Token,
		})

		if audios[i].Error {
			statuses[i].Status = statusError
		} else if audios[i].FinishedAt.Unix() <= 0 {
			statuses[i].Status = statusProcess
		} else {
			statuses[i].Status = statusDone
		}
	}

	service.ProcessOk(w, statuses)
}

// GetOne returns information about one audio
func (d *DecodeService) GetOne(w http.ResponseWriter, r *http.Request) {
	var req request.UserNameAudioToken
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	audio, err := as.GetByToken(req.Token)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUserAudio)
		return
	}
	status := request.AudioStatus{
		Name:  audio.Name,
		Date:  audio.PostedAt,
		Token: audio.Token,
	}

	if audio.Error {
		status.Status = statusError
	} else if audio.FinishedAt.Unix() <= 0 {
		status.Status = statusProcess
	} else {
		status.Status = statusDone
	}

	service.ProcessOk(w, status)
}

// Load returns audio by its token
func (d *DecodeService) Load(w http.ResponseWriter, r *http.Request) {
	var req request.UserNameAudioToken
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	audio, err := as.GetByToken(req.Token)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUserAudio)
		return
	}

	f, err := os.OpenFile(getLocation(req.Username, audio.Name), os.O_RDONLY, 0666)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUserAudio)
	}
	defer f.Close()

	io.Copy(w, f)
}

// Encode input audio
func (d *DecodeService) Encode(w http.ResponseWriter, r *http.Request) {
	var req request.UserAudio
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	_, err = storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	audio := &model.Audio{
		Name: req.AudioInit.Name,
		URL:  req.AudioInit.URL,
		Crypto: model.Crypto{
			Key: req.AudioInit.Key,
		},
		PostedAt:   time.Now(),
		FinishedAt: time.Time{},
		Token:      service.GenerateToken(),
	}
	as.Save(audio)

	path, _ := filepath.Abs(getBaseDir(req.Username))
	os.MkdirAll(path, 0777)
	go func(a storage.Audio, lAudio *model.Audio, path string) {
		rCh, eCh, data := startFfmpegEncrypt(req.AudioInit.URL, lAudio.Token+encryptExtension, req.AudioInit.Key, req.AudioInit.KID, path, getBaseSourceDir(req.Username))
		select {
		case res := <-rCh:
			if res.StatusCode != 0 {
				lAudio.Error = true
				as.Save(lAudio)
				io.Copy(os.Stdout, data)
				return
			}
			lAudio.Name = lAudio.Token + encryptExtension
			lAudio.FinishedAt = time.Now()
			as.Save(lAudio)
		case <-eCh:
			lAudio.Error = true
			as.Save(lAudio)
		}
	}(as, audio, path)

	service.ProcessOk(w, request.AudioToken{Token: audio.Token})
}

// Decode input audio
func (d *DecodeService) Decode(w http.ResponseWriter, r *http.Request) {
	var req request.UserAudio
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, service.ErrWrongFormat)
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	_, err = storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, service.ErrFindUser)
		return
	}

	audio := &model.Audio{
		Name: req.AudioInit.Name,
		URL:  req.AudioInit.URL,
		Crypto: model.Crypto{
			Key: req.AudioInit.Key,
		},
		PostedAt:   time.Now(),
		FinishedAt: time.Time{},
		Token:      service.GenerateToken(),
	}
	as.Save(audio)

	path, _ := filepath.Abs(getBaseDir(req.Username))
	os.MkdirAll(path, 0777)
	go func(a storage.Audio, lAudio *model.Audio, path string) {
		rCh, eCh, _ := startFfmpegDecrypt(req.AudioInit.URL, lAudio.Token+audioExtension, req.AudioInit.Key, path, getBaseSourceDir(req.Username))
		select {
		case res := <-rCh:
			if res.StatusCode != 0 {
				time.Sleep(4 * time.Second)
				lAudio.Error = true
				as.Save(lAudio)
				return
			}
			lAudio.Name = lAudio.Token + audioExtension
			lAudio.FinishedAt = time.Now()
			as.Save(lAudio)
		case <-eCh:
			lAudio.Error = true
			as.Save(lAudio)
		}
	}(as, audio, path)

	service.ProcessOk(w, request.AudioToken{Token: audio.Token})
}
