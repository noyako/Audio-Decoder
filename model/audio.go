package model

import "time"

// Audio represents encrypted/decrypted audio
type Audio struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	URL        string
	PostedAt   time.Time
	FinishedAt time.Time
	Token      string `gorm:"uniqueIndex"`
	Crypto     Crypto `gorm:"embedded"`
	Error      bool
}

// Crypto represents encryption/decryption parameters
type Crypto struct {
	Key                 string
	EncryptionType      string
	EncryptionDirection string
}
