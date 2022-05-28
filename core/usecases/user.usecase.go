package usecases

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/google/uuid"
)

type UserUsecase struct {
	repo   *repo.UserRepo
	mailer *services.Mailer
}

func NewUserUsecase(repo *repo.UserRepo, mailer *services.Mailer) *UserUsecase {
	return &UserUsecase{
		repo: repo, mailer: mailer,
	}
}

func (u *UserUsecase) GetFull(baseInfo dto.BaseUserInfo) (repo.User, error) {
	return u.repo.GetByEmail(baseInfo.Email)
}

func (u *UserUsecase) GetByID(id uuid.UUID) Result {
	user, err := u.repo.Get(id)
	if err != nil {
		return Result{nil, MakeUsecaseError("User not found.")}
	}
	return Result{user, nil}
}

func (u *UserUsecase) Create(createDto dto.CreateUser) Result {
	user, err := u.repo.Create(createDto)
	if err != nil {
		return Result{nil, MakeUsecaseError("User already exists.")}
	}

	go func() {
		to := []string{user.Email}
		err := u.mailer.SendVerifyEmail(to, user.ID.String())
		if err != nil {
			misc.ReportError("Cant sent verify email!")
		}
	}()

	return Result{repo.UserToBaseUser(user), nil}
}

func (u *UserUsecase) Delete(user dto.BaseUserInfo) Result {
	err := u.repo.Delete(user)
	if err != nil {
		return Result{false, MakeUsecaseError("Internal database error, try delete later.")}
	}
	return Result{true, nil}
}

func (u *UserUsecase) Update(id uuid.UUID, updateDto dto.UpdateUser) Result {
	user, err := u.repo.Update(id, updateDto)
	if err != nil {
		return Result{nil, MakeUsecaseError("Internal database error, try update later.")}
	}
	return Result{user, err}
}

func (u *UserUsecase) Verify(id uuid.UUID) Result {
	err := u.repo.Verify(id)
	if err != nil {
		return Result{false, MakeUsecaseError("Broken verify link, try later.")}
	}
	return Result{true, nil}
}
