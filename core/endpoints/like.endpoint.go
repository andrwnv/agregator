package endpoints

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/google/uuid"
)

type LikeEndpoint struct {
	likedRepo     *repo.LikedRepo
	userEndpoint  *UserEndpoint
	eventEndpoint *EventEndpoint
	placeEndpoint *PlaceEndpoint
}

func NewLikeEndpoint(
	likedRepo *repo.LikedRepo,
	userEndpoint *UserEndpoint,
	eventEndpoint *EventEndpoint,
	placeEndpoint *PlaceEndpoint) *LikeEndpoint {

	return &LikeEndpoint{
		likedRepo:     likedRepo,
		userEndpoint:  userEndpoint,
		eventEndpoint: eventEndpoint,
		placeEndpoint: placeEndpoint,
	}
}

func (e *LikeEndpoint) Get(userInfo dto.BaseUserInfo, page int, count int) Result {
	likes, err := e.likedRepo.Get(userInfo, page, count)
	if err != nil {
		return Result{nil, MakeEndpointError("Something went wrong, try later")}
	}
	var result []dto.LikedDto
	for _, like := range likes {
		result = append(result, repo.LikeToLike(like))
	}

	return Result{result, nil}
}

func (e *LikeEndpoint) Like(like dto.LikeDto, userInfo dto.BaseUserInfo) Result {
	user, err := e.userEndpoint.GetFull(userInfo)
	if err != nil {
		return Result{nil, MakeEndpointError("Cant find user for like something.")}
	}

	var liked repo.Liked
	var likeErr error
	if like.PlaceID != nil {
		place, getErr := e.placeEndpoint.GetFullPlace(*like.PlaceID)
		if getErr != nil {
			return Result{nil, MakeEndpointError("Cant find item to like")}
		}
		liked, likeErr = e.likedRepo.LikePlace(user, place)
	} else if like.EventID != nil {
		event, getErr := e.eventEndpoint.GetFullEvent(*like.EventID)
		if getErr != nil {
			return Result{nil, MakeEndpointError("Cant find item to like")}
		}
		liked, likeErr = e.likedRepo.LikeEvent(user, event)
	}

	if likeErr != nil {
		return Result{nil, MakeEndpointError("Already liked.")}
	}
	return Result{repo.LikeToLike(liked), nil}
}

func (e *LikeEndpoint) Dislike(id uuid.UUID, userInfo dto.BaseUserInfo) Result {
	user, err := e.userEndpoint.GetFull(userInfo)
	if err != nil {
		return Result{false, MakeEndpointError("Cant find user for dislike something.")}
	}

	if dislikeErr := e.likedRepo.Dislike(user, id); dislikeErr != nil {
		return Result{false, MakeEndpointError("Something went wrong, try later")}
	}
	return Result{true, nil}
}
