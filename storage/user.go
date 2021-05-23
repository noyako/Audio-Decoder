package storage

import "github.com/noyako/Audio-Decoder/model"

type User interface {
	GetByToken(string) (*model.User, error)
	GetByName(string) (*model.User, error)
	Save(*model.User) error
	Migrate()
}
