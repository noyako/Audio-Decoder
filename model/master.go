package model

type Master struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Database string
}
