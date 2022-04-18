package models

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/dto"
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

func DeleteUser(id string) (err error) {
	return core.ServerInst.Database.Where("id = ?", id).Delete(&User{}).Error
}

func CreateUser(u User) (err error) {
	return core.ServerInst.Database.Create(&u).Error
}

func GetByEmail(email string) (user User, err error) {
	err = core.ServerInst.Database.Where("email = ?", email).First(&user).Error
	return user, err
}

// ------------------ Conversations ------------------

func From(dto dto.CreateUser) User {
	passHash := sha1.New()
	passHash.Write([]byte(dto.Password))

	return User{
		ID:        uuid.New(),
		FirstName: dto.FirstName,
		LastName:  dto.SecondName,
		Email:     dto.Email,
		Password:  hex.EncodeToString(passHash.Sum(nil)),
	}
}

func To(user User) dto.BaseUserInfo {
	return dto.BaseUserInfo{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		SecondName: user.LastName,
		Email:      user.Email,
	}
}
