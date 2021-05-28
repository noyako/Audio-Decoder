package storage

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/noyako/Audio-Decoder/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// UserPostgres user for postgres database
type UserPostgres struct {
	db *gorm.DB
}

// NewUserPostgres UserPostgres constructor
func NewUserPostgres(db *sql.DB) (*UserPostgres, error) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	return &UserPostgres{
		db: gormDB,
	}, err
}

// GetByToken returns user by its token
func (u *UserPostgres) GetByToken(token string) (*model.User, error) {
	var user model.User
	u.db.Where("token = ?", token).First(&user)
	if (reflect.DeepEqual(user, model.User{})) {
		return &user, fmt.Errorf(errUserTokenNotFound, token)
	}
	return &user, nil
}

// GetByName returns user by its name
func (u *UserPostgres) GetByName(username string) (*model.User, error) {
	var user model.User
	u.db.Where("username = ?", username).First(&user)
	if (reflect.DeepEqual(user, model.User{})) {
		return &user, fmt.Errorf(errUserCredsNotFound, username)
	}
	return &user, nil
}

// Save user
func (u *UserPostgres) Save(user *model.User) error {
	result := u.db.Create(user)
	return result.Error
}

// Migrate database (user table only)
func (u *UserPostgres) Migrate() {
	u.db.AutoMigrate(&model.User{})
}
