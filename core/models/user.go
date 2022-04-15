package models

import (
	"github.com/andrwnv/event-aggregator/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID        uuid.UUID `gorm:"primaryKey"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
}

func CreateUser(u *User) (err error) {
	return core.ServerInst.Database.Create(u).Error
}

func GetByEmail(u *User, email string) (err error) {
	return core.ServerInst.Database.Where("email = ?", email).First(u).Error
}
