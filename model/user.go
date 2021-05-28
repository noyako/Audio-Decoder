package model

// User tenant information
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
	Token    string `gorm:"uniqueIndex"`
}
