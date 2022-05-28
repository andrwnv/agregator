package usecases

import (
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
)

type AuthUsecase struct {
	userUsecase *UserUsecase
}

func NewAuthUsecase(
	userUsecase *UserUsecase) *AuthUsecase {
	return &AuthUsecase{
		userUsecase: userUsecase,
	}
}

func (u *AuthUsecase) Login(credentials dto.LoginCredentials) Result {
	user, err := u.userUsecase.GetFull(dto.BaseUserInfo{
		Email: credentials.Email,
	})
	if err != nil {
		return Result{nil, err}
	}

	authSuccess := services.Login(credentials, services.LoginInfo{
		Email:    user.Email,
		Password: user.Password,
	})
	if authSuccess {
		token := core.SERVER.JwtService.GenerateToken(credentials.Email, repo.UserToBaseUser(user))
		if token != "" {
			return Result{token, nil}
		}
	}

	return Result{nil, MakeUsecaseError("Invalid credentials.")}
}
