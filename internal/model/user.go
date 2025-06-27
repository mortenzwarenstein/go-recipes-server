package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Base
	Email        string `gorm:"size:255;not null;unique" json:"email"`
	Password     string `gorm:"size:255;not null" json:"-"`
	RefreshToken string `gorm:"size:512;" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if err := u.Base.BeforeCreate(tx); err != nil {
		return err
	}

	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}

	return nil
}
