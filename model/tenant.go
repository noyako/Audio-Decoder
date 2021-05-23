package model

import "time"

type Tenant struct {
	ID        uint `gorm:"primaryKey"`
	Published time.Time
	URL       string
	Status    string
	ResultURL string
}
