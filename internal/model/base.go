package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        string         `gorm:"primary_key;type:uuid" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return
}
