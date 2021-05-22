package storage

import (
	"database/sql"
	"fmt"

	"github.com/noyako/Audio-Decoder/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AudioPostgres struct {
	db *gorm.DB
}

func NewAudioPostgres(db *sql.DB) (*AudioPostgres, error) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	return &AudioPostgres{
		db: gormDB,
	}, err
}

func (a *AudioPostgres) GetByToken(token string) (*model.Audio, error) {
	var audio model.Audio
	a.db.Where("token = ?", token).First(&audio)
	if (audio == model.Audio{}) {
		return &audio, fmt.Errorf(errAudioToknenNotFound, token)
	}
	return &audio, nil
}

func (a *AudioPostgres) GetByOwner(user *model.User) ([]*model.Audio, error) {
	var audio []*model.Audio
	a.db.Where("owner_id = ?", user.ID).Find(&audio)
	if audio == nil {
		return audio, fmt.Errorf(errAudioUserNotFound, user.ID)
	}
	return audio, nil
}

func (a *AudioPostgres) Save(audio *model.Audio) error {
	result := a.db.Create(audio)
	return result.Error
}
