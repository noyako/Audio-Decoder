package storage

import (
	"github.com/noyako/Audio-Decoder/model"
)

type Audio interface {
	GetByToken(string) (*model.Audio, error)
	GetByOwner(*model.User) ([]*model.Audio, error)
	Save(*model.Audio) error
}
