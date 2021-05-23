package storage

import (
	"github.com/noyako/Audio-Decoder/model"
)

type Audio interface {
	GetByToken(string) (*model.Audio, error)
	GetAll() ([]*model.Audio, error)
	Save(*model.Audio) error
	Migrate()
}
