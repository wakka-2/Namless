package models

import (
	"time"

	"gorm.io/gorm"
)

type Data struct {
	ID        string `json:"id,omitempty" gorm:"primary_key"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
