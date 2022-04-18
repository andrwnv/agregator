package repo

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

// ----- UserRepo methods

type UserRepoCrud interface {
	Create(dto dto.CreateUser) error
	Delete(dto dto.BaseUserInfo) error
	GetByEmail(email string) (User, error)
	// Update(dto ...) error
}

type UserRepo struct {
	Repo *PgRepo
}

func NewUserRepo(repo *PgRepo) *UserRepo {
	return &UserRepo{
		Repo: repo,
	}
}

func (repo *UserRepo) Create(dto dto.CreateUser) error {
	passHash := sha1.New()
	passHash.Write([]byte(dto.Password))

	return repo.Repo.Database.Create(&User{
		ID:        uuid.New(),
		FirstName: dto.FirstName,
		LastName:  dto.SecondName,
		Email:     dto.Email,
		Password:  hex.EncodeToString(passHash.Sum(nil)),
	}).Error
}

func (repo *UserRepo) Delete(dto dto.BaseUserInfo) error {
	return repo.Repo.Database.Exec("DELETE FROM users WHERE id = ?", dto.ID).Error
}

func (repo *UserRepo) GetByEmail(email string) (user User, err error) {
	err = core.ServerInst.Database.Where("email = ?", email).First(&user).Error
	return user, err
}
