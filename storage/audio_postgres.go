package storage

import (
	"database/sql"
	"fmt"

	"github.com/noyako/Audio-Decoder/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AudioPostgres audio for postgres database
type AudioPostgres struct {
	db *gorm.DB
}

// NewAudioPostgres AudioPostgres constructor
func NewAudioPostgres(db *sql.DB) (*AudioPostgres, error) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	return &AudioPostgres{
		db: gormDB,
	}, err
}

// GetByToken returns audio by its token
func (a *AudioPostgres) GetByToken(token string) (*model.Audio, error) {
	var audio model.Audio
	a.db.Where("token = ?", token).First(&audio)
	if (audio == model.Audio{}) {
		return &audio, fmt.Errorf(errAudioToknenNotFound, token)
	}
	return &audio, nil
}

// GetAll returns all audios
func (a *AudioPostgres) GetAll() ([]*model.Audio, error) {
	var audio []*model.Audio
	a.db.Find(&audio)
	if audio == nil {
		return audio, fmt.Errorf(errAudioNotFound)
	}
	return audio, nil
}

// Save audio
func (a *AudioPostgres) Save(audio *model.Audio) error {
	result := a.db.Save(audio)
	return result.Error
}

// Remove audio
func (a *AudioPostgres) Remove(audio *model.Audio) error {
	result := a.db.Delete(audio)
	return result.Error
}

// Migrate database tables
func (a *AudioPostgres) Migrate() {
	a.db.AutoMigrate(&model.Audio{})
}
