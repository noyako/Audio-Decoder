package storage

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/noyako/Audio-Decoder/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserPostgres struct {
	db *gorm.DB
}

func NewUserPostgres(db *sql.DB) (*AudioPostgres, error) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	return &AudioPostgres{
		db: gormDB,
	}, err
}

func (u *UserPostgres) GetByToken(token string) (*model.User, error) {
	var user model.User
	u.db.Where("token = ?", token).First(&user)
	if (reflect.DeepEqual(user, model.User{})) {
		return &user, fmt.Errorf(errUserTokenNotFound, token)
	}
	return &user, nil
}

func (u *UserPostgres) GetByCreds(username, password string) (*model.User, error) {
	var user model.User
	u.db.Where("username = ? AND password = ?", username, password).First(&user)
	if (reflect.DeepEqual(user, model.User{})) {
		return &user, fmt.Errorf(errUserCredsNotFound, username)
	}
	return &user, nil
}

func (u *UserPostgres) Save(user *model.User) error {
	result := u.db.Create(user)
	return result.Error
}
