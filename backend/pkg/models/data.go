/*
Package models offers DB related structures.
*/
package models

import (
	"time"

	"gorm.io/gorm"
)

// Data models a (key, value) data item.
type Data struct {
	ID        string `json:"id,omitempty" gorm:"primary_key"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
