package model

// Master table for tenants
type Master struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Database string
}
