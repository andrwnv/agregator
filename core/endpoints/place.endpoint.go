package endpoints

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/google/uuid"
)

type PlaceEndpoint struct {
	placeRepo    *repo.PlaceRepo
	userEndpoint *UserEndpoint
	regionRepo   *repo.RegionRepo
}

func NewPlaceEndpoint(
	placeRepo *repo.PlaceRepo,
	userEndpoint *UserEndpoint,
	regionRepo *repo.RegionRepo) *PlaceEndpoint {
	return &PlaceEndpoint{
		placeRepo:    placeRepo,
		userEndpoint: userEndpoint,
		regionRepo:   regionRepo,
	}
}

func (e *PlaceEndpoint) Get(id uuid.UUID) Result {
	place, err := e.placeRepo.Get(id)
	placePhotos, _ := e.placeRepo.GetImages(id)
	return Result{repo.PlaceToPlace(place, placePhotos), err}
}

func (e *PlaceEndpoint) GetFullPlace(id uuid.UUID) (repo.Place, error) {
	return e.placeRepo.Get(id)
}

func (e *PlaceEndpoint) GetFullPlaceComment(id uuid.UUID) (repo.PlaceComment, error) {
	return e.placeRepo.GetCommentByID(id)
}

func (e *PlaceEndpoint) Create(createDto dto.CreatePlace, userInfo dto.BaseUserInfo) Result {
	user, err := e.userEndpoint.GetFull(userInfo)
	if err != nil {
		return Result{nil, MakeEndpointError("Cant find user for create place.")}
	}

	region, err := e.regionRepo.GetByRegionID(createDto.RegionID)
	if err != nil {
		return Result{nil, MakeEndpointError("Cant find selected country.")}
	}

	place, err := e.placeRepo.Create(createDto, user, region)
	return Result{repo.PlaceToPlace(place, []string{}), err}
}

func (e *PlaceEndpoint) Update(id uuid.UUID, updateDto dto.UpdatePlace, userInfo dto.BaseUserInfo) Result {
	place, err := e.placeRepo.Get(id)
	if err != nil {
		return Result{nil, err}
	}

	if userInfo.ID != place.CreatedBy.ID.String() {
		return Result{nil, MakeEndpointError("Isn't your place!")}
	}

	place.Region, err = e.regionRepo.GetByRegionID(updateDto.RegionID)
	if err != nil {
		return Result{nil, MakeEndpointError("Cant find selected country.")}
	}

	err = e.placeRepo.Update(place.ID, updateDto, place.Region)
	return Result{err != nil, err}
}

func (e *PlaceEndpoint) Delete(id uuid.UUID, userInfo dto.BaseUserInfo) Result {
	place, err := e.placeRepo.Get(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != place.CreatedBy.ID.String() {
		return Result{false, MakeEndpointError("Isn't your place!")}
	}

	err = e.placeRepo.Delete(id)
	return Result{err != nil, err}
}

// ----- PlaceEndpoint: Images -----

func (e *PlaceEndpoint) UpdatePlaceImages(id uuid.UUID, userInfo dto.BaseUserInfo,
	filesToCreate []string, filesToDelete []string) Result {

	place, err := e.placeRepo.Get(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != place.CreatedBy.ID.String() {
		return Result{false, MakeEndpointError("Isn't your place!")}
	}

	for _, url := range filesToCreate {
		err := e.placeRepo.CreateImages(place.ID, url)
		if err != nil {
			return Result{false, err}
		}
	}

	for _, url := range filesToDelete {
		err := e.placeRepo.DeleteImages(url)
		// TODO: delete photos from dir.
		if err != nil {
			return Result{false, err}
		}
	}

	return Result{true, nil}
}

// ----- PlaceEndpoint: Comments -----

func (e *PlaceEndpoint) CreateComment(createDto dto.CreatePlaceCommentDto, userInfo dto.BaseUserInfo) Result {
	user, err := e.userEndpoint.GetFull(userInfo)
	if err != nil {
		return Result{nil, err}
	}
	place, err := e.placeRepo.Get(uuid.MustParse(createDto.LinkedPlaceID))
	if err != nil {
		return Result{false, err}
	}

	comment, err := e.placeRepo.CreateComment(createDto, user, place)
	if err != nil {
		return Result{nil, MakeEndpointError("Failed to create comment.")}
	}

	return Result{repo.CommentToCommentDto(comment), nil}
}

func (e *PlaceEndpoint) GetComments(placeId uuid.UUID, page int, count int) Result {
	comments, err := e.placeRepo.GetComments(placeId, page, count)
	if err != nil {
		return Result{nil, MakeEndpointError("Failed to create comment.")}
	}

	var result []dto.PlaceCommentDto
	for _, value := range comments {
		result = append(result, repo.CommentToCommentDto(value))
	}

	return Result{result, nil}
}

func (e *PlaceEndpoint) DeleteComment(commentId uuid.UUID, userInfo dto.BaseUserInfo) Result {
	comment, err := e.placeRepo.GetCommentByID(commentId)
	if err != nil {
		return Result{nil, err}
	}
	if userInfo.ID != comment.CreatedBy.ID.String() {
		return Result{nil, MakeEndpointError("Isn't your comment!")}
	}

	err = e.placeRepo.DeleteComments(commentId)
	if err != nil {
		return Result{false, MakeEndpointError("Cant delete comment(s).")}
	}
	return Result{true, nil}
}

func (e *PlaceEndpoint) UpdateComment(id uuid.UUID, updateDto dto.UpdatePlaceCommentDto, userInfo dto.BaseUserInfo) Result {
	comment, err := e.placeRepo.GetCommentByID(id)
	if err != nil {
		return Result{nil, err}
	}
	if userInfo.ID != comment.CreatedBy.ID.String() {
		return Result{nil, MakeEndpointError("Isn't your comment!")}
	}

	err = e.placeRepo.UpdateComment(id, updateDto)
	if err != nil {
		return Result{false, MakeEndpointError("Cant update comment(s).")}
	}
	return Result{true, nil}
}
