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

func (a *AudioPostgres) GetAll() ([]*model.Audio, error) {
	var audio []*model.Audio
	a.db.Find(&audio)
	if audio == nil {
		return audio, fmt.Errorf(errAudioNotFound)
	}
	return audio, nil
}

func (a *AudioPostgres) Save(audio *model.Audio) error {
	result := a.db.Save(audio)
	return result.Error
}

func (u *AudioPostgres) Migrate() {
	u.db.AutoMigrate(&model.Audio{})
}
