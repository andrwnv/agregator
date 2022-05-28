package usecases

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/google/uuid"
)

type LikeUsecase struct {
	likedRepo    *repo.LikedRepo
	userUsecase  *UserUsecase
	eventUsecase *EventUsecase
	placeUsecase *PlaceUsecase
}

func NewLikeUsecase(
	likedRepo *repo.LikedRepo,
	userUsecase *UserUsecase,
	eventUsecase *EventUsecase,
	placeUsecase *PlaceUsecase) *LikeUsecase {

	return &LikeUsecase{
		likedRepo:    likedRepo,
		userUsecase:  userUsecase,
		eventUsecase: eventUsecase,
		placeUsecase: placeUsecase,
	}
}

func (u *LikeUsecase) Get(userInfo dto.BaseUserInfo, page int, count int) Result {
	likes, err := u.likedRepo.Get(userInfo, page, count)
	if err != nil {
		return Result{nil, MakeUsecaseError("Something went wrong, try later")}
	}
	var result []dto.LikedDto
	for _, like := range likes {
		result = append(result, repo.LikeToLike(like))
	}

	return Result{result, nil}
}

func (u *LikeUsecase) Like(like dto.LikeDto, userInfo dto.BaseUserInfo) Result {
	user, err := u.userUsecase.GetFull(userInfo)
	if err != nil {
		return Result{nil, MakeUsecaseError("Cant find user for like something.")}
	}

	var liked repo.Liked
	var likeErr error
	if like.PlaceID != nil {
		place, getErr := u.placeUsecase.GetFullPlace(*like.PlaceID)
		if getErr != nil {
			return Result{nil, MakeUsecaseError("Cant find item to like")}
		}
		liked, likeErr = u.likedRepo.LikePlace(user, place)
	} else if like.EventID != nil {
		event, getErr := u.eventUsecase.GetFullEvent(*like.EventID)
		if getErr != nil {
			return Result{nil, MakeUsecaseError("Cant find item to like")}
		}
		liked, likeErr = u.likedRepo.LikeEvent(user, event)
	}

	if likeErr != nil {
		return Result{nil, MakeUsecaseError("Already liked.")}
	}
	return Result{repo.LikeToLike(liked), nil}
}

func (u *LikeUsecase) Dislike(id uuid.UUID, userInfo dto.BaseUserInfo) Result {
	user, err := u.userUsecase.GetFull(userInfo)
	if err != nil {
		return Result{false, MakeUsecaseError("Cant find user for dislike something.")}
	}

	if dislikeErr := u.likedRepo.Dislike(user, id); dislikeErr != nil {
		return Result{false, MakeUsecaseError("Something went wrong, try later")}
	}
	return Result{true, nil}
}
