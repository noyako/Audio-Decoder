package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/noyako/Audio-Decoder/model"
	"github.com/noyako/Audio-Decoder/request"
	"github.com/noyako/Audio-Decoder/service"
	"github.com/noyako/Audio-Decoder/storage"
)

func (d *DecodeService) New(w http.ResponseWriter, r *http.Request) {
	var req request.UserName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	db, err := d.adb.Create(req.Username)
	if err != nil {
		service.ProcessServerError(w, "Cannot create user")
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, "Cannot create user")
		return
	}

	as.Migrate()

	// path, _ := os.Getwd()
	os.MkdirAll(path.Join("files", req.Username), 0777)
}

func (d *DecodeService) GetAll(w http.ResponseWriter, r *http.Request) {
	var req request.UserName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}
	as.Migrate()

	audios, err := as.GetAll()
	if err != nil {
		service.ProcessServerError(w, "Cannot get user audios")
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
			statuses[i].Status = "error"
		} else if audios[i].FinishedAt.Unix() <= 0 {
			statuses[i].Status = "processing"
		} else {
			statuses[i].Status = "done"
		}
	}

	service.ProcessOk(w, statuses)
}

func (d *DecodeService) GetOne(w http.ResponseWriter, r *http.Request) {
	var req request.UserNameAudioToken
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}

	audio, err := as.GetByToken(req.Token)
	if err != nil {
		service.ProcessServerError(w, "Cannot find audio")
		return
	}
	status := request.AudioStatus{
		Name:  audio.Name,
		Date:  audio.PostedAt,
		Token: audio.Token,
	}

	if audio.Error {
		status.Status = "error"
	} else if audio.FinishedAt.Unix() <= 0 {
		status.Status = "processing"
	} else {
		status.Status = "done"
	}

	service.ProcessOk(w, status)
}

func (d *DecodeService) Load(w http.ResponseWriter, r *http.Request) {
	var req request.UserNameAudioToken
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}

	audio, err := as.GetByToken(req.Token)
	if err != nil {
		service.ProcessServerError(w, "Cannot find audio")
		return
	}

	f, err := os.OpenFile(getLocation(req.Username, audio.Name), os.O_RDONLY, 0666)
	if err != nil {
		service.ProcessServerError(w, "Cannot find audio")
	}
	defer f.Close()

	io.Copy(w, f)
}

func (d *DecodeService) Encode(w http.ResponseWriter, r *http.Request) {
	var req request.UserAudio
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}

	_, err = storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
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

	dir, _ := os.Getwd()
	go func(a storage.Audio, lAudio *model.Audio, path string) {
		rCh, eCh, _ := startFfmpegEncrypt(req.AudioInit.URL, lAudio.Token+".crypt", req.AudioInit.Key, req.AudioInit.KID, path)
		select {
		case res := <-rCh:
			if res.StatusCode != 0 {
				lAudio.Error = true
				as.Save(lAudio)
				return
			}
			lAudio.Name = lAudio.Token + ".crypt"
			lAudio.FinishedAt = time.Now()
			as.Save(lAudio)
		case <-eCh:
			lAudio.Error = true
			as.Save(lAudio)
		}
	}(as, audio, path.Join(dir, "files", req.Username))

	service.ProcessOk(w, request.AudioToken{audio.Token})
}

func (d *DecodeService) Decode(w http.ResponseWriter, r *http.Request) {
	var req request.UserAudio
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.ProcessBadFormat(w, "Request json in wrong format")
		return
	}

	db, err := d.adb.Get(req.Username)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}

	_, err = storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
		return
	}

	as, err := storage.NewAudioPostgres(db)
	if err != nil {
		service.ProcessServerError(w, "Cannot get user")
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

	dir, _ := os.Getwd()
	go func(a storage.Audio, lAudio *model.Audio, path string) {
		rCh, eCh, _ := startFfmpegDecrypt(req.AudioInit.URL, lAudio.Token+".mp3", req.AudioInit.Key, path)
		select {
		case res := <-rCh:
			if res.StatusCode != 0 {
				lAudio.Error = true
				as.Save(lAudio)
				return
			}
			lAudio.Name = lAudio.Token + ".mp3"
			lAudio.FinishedAt = time.Now()
			as.Save(lAudio)
		case <-eCh:
			lAudio.Error = true
			as.Save(lAudio)
		}
	}(as, audio, path.Join(dir, "files", req.Username))
	// data, _ := ioutil.ReadAll(output)
	// fmt.Println(string(data))

	// res, _, output = startFfmpegDecrypt("output.crypt", "origin.mp3", "2BB80D537B1DA3E38BD30361AA855686")

	// fmt.Println(<-res)
	// data, _ = ioutil.ReadAll(output)
	// fmt.Println(string(data))

	service.ProcessOk(w, request.AudioToken{audio.Token})
}
