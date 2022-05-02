package endpoints

import (
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
)

type AuthEndpoint struct {
	userEndpoint *UserEndpoint
}

func NewAuthEndpoint(
	userEndpoint *UserEndpoint) *AuthEndpoint {
	return &AuthEndpoint{
		userEndpoint: userEndpoint,
	}
}

func (e *AuthEndpoint) Login(credentials dto.LoginCredentials) Result {
	user, err := e.userEndpoint.GetFull(dto.BaseUserInfo{
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

	return Result{nil, MakeEndpointError("Invalid credentials.")}
}
