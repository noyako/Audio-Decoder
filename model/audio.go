package model

import "time"

type Audio struct {
	ID         uint `gorm:"primaryKey"`
	OwnerID    uint
	Name       string
	Location   string
	PostedAt   time.Time
	FinishedAt time.Time
	Token      string `gorm:"uniqueIndex"`
	Crypto     Crypto `gorm:"embedded"`
}

type Crypto struct {
	Key                 string
	EncryptionType      string
	EncryptionDirection string
}
