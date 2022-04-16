package services

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/andrwnv/event-aggregator/core/dto"
)

type AuthService interface {
	Login(email string, password string) bool
}

type LoginInfo struct {
	Email    string
	Password string
}

func Login(dto dto.LoginCredentials, i LoginInfo) bool {
	passHash := sha1.New()
	passHash.Write([]byte(dto.Password))

	return i.Email == dto.Email && hex.EncodeToString(passHash.Sum(nil)) == i.Password
}
