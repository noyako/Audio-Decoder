package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/noyako/Audio-Decoder/service"
	"github.com/noyako/Audio-Decoder/storage"
)

func (d *DecodeService) New(w http.ResponseWriter, r *http.Request) {
	var req UserName
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
}

func (d *DecodeService) GetAll(w http.ResponseWriter, r *http.Request) {
	var req UserName
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

	audios, err := as.GetAll()
	if err != nil {
		service.ProcessServerError(w, "Cannot get user audios")
		return
	}

	var statuses []Status
	for i := range audios {
		statuses = append(statuses, Status{
			Name:  audios[i].Name,
			Date:  audios[i].PostedAt,
			Token: audios[i].Token,
		})

		if (audios[i].FinishedAt == time.Time{}) {
			statuses[i].Status = "processing"
		} else {
			statuses[i].Status = "done"
			statuses[i].Token = audios[i].Token
		}
	}

	service.ProcessOk(w, statuses)
}

func (d *DecodeService) GetOne(w http.ResponseWriter, r *http.Request) {
	var req Audio
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
	status := Status{
		Name:  audio.Name,
		Date:  audio.PostedAt,
		Token: audio.Token,
	}

	if (audio.FinishedAt == time.Time{}) {
		status.Status = "processing"
	} else {
		status.Status = "done"
		status.Token = audio.Token
	}

	service.ProcessOk(w, status)
}

func (d *DecodeService) Load(w http.ResponseWriter, r *http.Request) {
	var req Audio
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
