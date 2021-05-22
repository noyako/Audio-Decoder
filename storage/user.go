package storage

import "github.com/noyako/Audio-Decoder/model"

type User interface {
	GetByToken(string) (*model.User, error)
	GetByCreds(string, string) (*model.User, error)
	Save(*model.User) error
}
