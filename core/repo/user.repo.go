package repo

import (
	"crypto/sha1"
	"encoding/hex"
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
	BirthDay  int       `gorm:"null,type:bigint"`
	Verified  bool      `gorm:"not null;default: false"`
	PhotoUrl  string    `gorm:"null"`
}

// ----- UserRepo methods -----

type UserRepoCrud interface {
	Create(dto dto.CreateUser) (user User, err error)
	Delete(dto dto.BaseUserInfo) error
	GetByEmail(email string) (User, error)
	Update(uuid uuid.UUID, dto dto.UpdateUser) (dto.BaseUserInfo, error)
	Verify(uuid uuid.UUID) error
}

type UserRepo struct {
	repo *PgRepo
}

func NewUserRepo(repo *PgRepo) *UserRepo {
	_ = repo.Database.AutoMigrate(&User{})

	return &UserRepo{
		repo: repo,
	}
}

func (repo *UserRepo) Create(dto dto.CreateUser) (user User, err error) {
	passHash := sha1.New()
	passHash.Write([]byte(dto.Password))

	user = User{
		ID:        uuid.New(),
		FirstName: dto.FirstName,
		LastName:  dto.SecondName,
		Email:     dto.Email,
		Password:  hex.EncodeToString(passHash.Sum(nil)),
	}

	return user, repo.repo.Database.Create(&user).Error
}

func (repo *UserRepo) Delete(dto dto.BaseUserInfo) error {
	return repo.repo.Database.Exec("DELETE FROM users WHERE id = ?", dto.ID).Error
}

func (repo *UserRepo) GetByEmail(email string) (user User, err error) {
	err = repo.repo.Database.Where("email = ?", email).First(&user).Error
	return user, err
}

func (repo *UserRepo) Update(uuid uuid.UUID, dto dto.UpdateUser) (dto.BaseUserInfo, error) {
	var user User
	repo.repo.Database.Where("id = ?", uuid).First(&user)

	user.FirstName = dto.FirstName
	user.LastName = dto.SecondName
	user.BirthDay = dto.BirthDay
	user.Password = dto.Password
	user.PhotoUrl = dto.PhotoUrl

	return To(user), repo.repo.Database.Save(&user).Error
}

func (repo *UserRepo) Verify(uuid uuid.UUID) error {
	var user User
	repo.repo.Database.Where("id = ?", uuid).First(&user)
	user.Verified = true

	return repo.repo.Database.Save(&user).Error
}

// ----- Conversations -----

func To(user User) dto.BaseUserInfo {
	return dto.BaseUserInfo{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		SecondName: user.LastName,
		Email:      user.Email,
		BirthDay:   user.BirthDay,
		Verified:   user.Verified,
		PhotoUrl:   user.PhotoUrl,
	}
}

func ToUpdateDto(user User) dto.UpdateUser {
	return dto.UpdateUser{
		FirstName:  user.FirstName,
		SecondName: user.LastName,
		BirthDay:   user.BirthDay,
		Password:   user.Password,
		PhotoUrl:   user.PhotoUrl,
	}
}
