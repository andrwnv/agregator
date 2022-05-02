package endpoints

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/google/uuid"
)

type UserEndpoint struct {
	repo   *repo.UserRepo
	mailer *services.Mailer
}

func NewUserEndpoint(repo *repo.UserRepo, mailer *services.Mailer) *UserEndpoint {
	return &UserEndpoint{
		repo: repo, mailer: mailer,
	}
}

func (e *UserEndpoint) GetFull(baseInfo dto.BaseUserInfo) (repo.User, error) {
	return e.repo.GetByEmail(baseInfo.Email)
}

func (e *UserEndpoint) GetByID(id uuid.UUID) Result {
	user, err := e.repo.Get(id)
	if err != nil {
		return Result{nil, MakeEndpointError("User not found.")}
	}
	return Result{user, nil}
}

func (e *UserEndpoint) Create(createDto dto.CreateUser) Result {
	user, err := e.repo.Create(createDto)
	if err != nil {
		return Result{nil, MakeEndpointError("User already exists.")}
	}

	go func() {
		to := []string{user.Email}
		err := e.mailer.SendVerifyEmail(to, user.ID.String())
		if err != nil {
			misc.ReportError("Cant sent verify email!")
		}
	}()

	return Result{repo.UserToBaseUser(user), nil}
}

func (e *UserEndpoint) Delete(user dto.BaseUserInfo) Result {
	err := e.repo.Delete(user)
	if err != nil {
		return Result{false, MakeEndpointError("Internal database error, try delete later.")}
	}
	return Result{true, nil}
}

func (e *UserEndpoint) Update(id uuid.UUID, updateDto dto.UpdateUser) Result {
	user, err := e.repo.Update(id, updateDto)
	if err != nil {
		return Result{nil, MakeEndpointError("Internal database error, try update later.")}
	}
	return Result{user, err}
}

func (e *UserEndpoint) Verify(id uuid.UUID) Result {
	err := e.repo.Verify(id)
	if err != nil {
		return Result{false, MakeEndpointError("Broken verify link, try later.")}
	}
	return Result{true, nil}
}
