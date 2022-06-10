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

func (u *User) BeforeDelete(tx *gorm.DB) error {
	var events []Event
	if err := tx.Table("events").Where("created_by_id = ?", u.ID).Find(&events).Unscoped().Delete(&events).Error; err != nil {
		return err
	}

	var places []Place
	if err := tx.Table("places").Where("created_by_id = ?", u.ID).Find(&places).Unscoped().Delete(&places).Error; err != nil {
		return err
	}

	var stories []UserStory
	if err := tx.Table("user_stories").Where("created_by_id = ?", u.ID).Find(&stories).Unscoped().Delete(&stories).Error; err != nil {
		return err
	}

	var likes []Liked
	if err := tx.Table("likeds").Where("user_id = ?", u.ID).Find(&likes).Unscoped().Delete(&likes).Error; err != nil {
		return err
	}

	return nil
}

// ----- UserRepo methods -----

type UserRepoCrud interface {
	Create(dto dto.CreateUser) (user User, err error)
	Get(uuid uuid.UUID) (dto.BaseUserInfo, error)
	GetByEmail(email string) (User, error)
	Delete(dto dto.BaseUserInfo) error
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
	return repo.repo.Database.Delete(&User{ID: uuid.MustParse(dto.ID)}).Error
}

func (repo *UserRepo) GetByEmail(email string) (user User, err error) {
	err = repo.repo.Database.Where("email = ?", email).First(&user).Error
	return user, err
}

func (repo *UserRepo) Get(uuid uuid.UUID) (u dto.BaseUserInfo, err error) {
	var user User
	err = repo.repo.Database.Where("id = ?", uuid).First(&user).Error
	return UserToBaseUser(user), err
}

func (repo *UserRepo) Update(uuid uuid.UUID, dto dto.UpdateUser) (dto.BaseUserInfo, error) {
	var user User
	repo.repo.Database.Where("id = ?", uuid).First(&user)

	user.FirstName = dto.FirstName
	user.LastName = dto.SecondName
	user.BirthDay = dto.BirthDay
	user.Password = dto.Password
	user.PhotoUrl = dto.PhotoUrl

	return UserToBaseUser(user), repo.repo.Database.Save(&user).Error
}

func (repo *UserRepo) Verify(uuid uuid.UUID) error {
	var user User
	repo.repo.Database.Where("id = ?", uuid).First(&user)
	user.Verified = true

	return repo.repo.Database.Save(&user).Error
}

// ----- Conversations -----

func UserToBaseUser(user User) dto.BaseUserInfo {
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

func UserToUpdateUserDto(user User) dto.UpdateUser {
	return dto.UpdateUser{
		FirstName:  user.FirstName,
		SecondName: user.LastName,
		BirthDay:   user.BirthDay,
		Password:   user.Password,
		PhotoUrl:   user.PhotoUrl,
	}
}
