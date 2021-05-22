package model

type User struct {
	ID             uint `gorm:"primaryKey"`
	Username       string
	Password       string
	Token          string  `gorm:"uniqueIndex"`
	PublishedAudio []Audio `gorm:"foreignKey:OwnerID"`
}
