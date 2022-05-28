package usecases

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/google/uuid"
)

type PlaceUsecase struct {
	placeRepo   *repo.PlaceRepo
	userUsecase *UserUsecase
	regionRepo  *repo.RegionRepo
}

func NewPlaceUsecase(
	placeRepo *repo.PlaceRepo,
	userUsecase *UserUsecase,
	regionRepo *repo.RegionRepo) *PlaceUsecase {
	return &PlaceUsecase{
		placeRepo:   placeRepo,
		userUsecase: userUsecase,
		regionRepo:  regionRepo,
	}
}

func (u *PlaceUsecase) Get(id uuid.UUID) Result {
	place, err := u.placeRepo.Get(id)
	placePhotos, _ := u.placeRepo.GetImages(id)
	return Result{repo.PlaceToPlace(place, placePhotos), err}
}

func (u *PlaceUsecase) GetFullPlace(id uuid.UUID) (repo.Place, error) {
	return u.placeRepo.Get(id)
}

func (u *PlaceUsecase) GetFullPlaceComment(id uuid.UUID) (repo.PlaceComment, error) {
	return u.placeRepo.GetCommentByID(id)
}

func (u *PlaceUsecase) Create(createDto dto.CreatePlace, userInfo dto.BaseUserInfo) Result {
	user, err := u.userUsecase.GetFull(userInfo)
	if err != nil {
		return Result{nil, MakeUsecaseError("Cant find user for create place.")}
	}

	region, err := u.regionRepo.GetByRegionID(createDto.RegionID)
	if err != nil {
		return Result{nil, MakeUsecaseError("Cant find selected country.")}
	}

	place, err := u.placeRepo.Create(createDto, user, region)
	return Result{repo.PlaceToPlace(place, []string{}), err}
}

func (u *PlaceUsecase) Update(id uuid.UUID, updateDto dto.UpdatePlace, userInfo dto.BaseUserInfo) Result {
	place, err := u.placeRepo.Get(id)
	if err != nil {
		return Result{nil, err}
	}

	if userInfo.ID != place.CreatedBy.ID.String() {
		return Result{nil, MakeUsecaseError("Isn't your place!")}
	}

	place.Region, err = u.regionRepo.GetByRegionID(updateDto.RegionID)
	if err != nil {
		return Result{nil, MakeUsecaseError("Cant find selected country.")}
	}

	err = u.placeRepo.Update(place.ID, updateDto, place.Region)
	return Result{err != nil, err}
}

func (u *PlaceUsecase) Delete(id uuid.UUID, userInfo dto.BaseUserInfo) Result {
	place, err := u.placeRepo.Get(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != place.CreatedBy.ID.String() {
		return Result{false, MakeUsecaseError("Isn't your place!")}
	}

	err = u.placeRepo.Delete(id)
	return Result{err != nil, err}
}

// ----- PlaceUsecase: Images -----

func (u *PlaceUsecase) UpdatePlaceImages(id uuid.UUID, userInfo dto.BaseUserInfo,
	filesToCreate []string, filesToDelete []string) Result {

	place, err := u.placeRepo.Get(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != place.CreatedBy.ID.String() {
		return Result{false, MakeUsecaseError("Isn't your place!")}
	}

	for _, url := range filesToCreate {
		err := u.placeRepo.CreateImages(place.ID, url)
		if err != nil {
			return Result{false, err}
		}
	}

	for _, url := range filesToDelete {
		err := u.placeRepo.DeleteImages(url)
		// TODO: delete photos from dir.
		if err != nil {
			return Result{false, err}
		}
	}

	return Result{true, nil}
}

// ----- PlaceUsecase: Comments -----

func (u *PlaceUsecase) CreateComment(createDto dto.CreatePlaceCommentDto, userInfo dto.BaseUserInfo) Result {
	user, err := u.userUsecase.GetFull(userInfo)
	if err != nil {
		return Result{nil, err}
	}
	place, err := u.placeRepo.Get(uuid.MustParse(createDto.LinkedPlaceID))
	if err != nil {
		return Result{false, err}
	}

	comment, err := u.placeRepo.CreateComment(createDto, user, place)
	if err != nil {
		return Result{nil, MakeUsecaseError("Failed to create comment.")}
	}

	return Result{repo.CommentToCommentDto(comment), nil}
}

func (u *PlaceUsecase) GetComments(placeId uuid.UUID, page int, count int) Result {
	comments, err := u.placeRepo.GetComments(placeId, page, count)
	if err != nil {
		return Result{nil, MakeUsecaseError("Failed to create comment.")}
	}

	var result []dto.PlaceCommentDto
	for _, value := range comments {
		result = append(result, repo.CommentToCommentDto(value))
	}

	return Result{result, nil}
}

func (u *PlaceUsecase) DeleteComment(commentId uuid.UUID, userInfo dto.BaseUserInfo) Result {
	comment, err := u.placeRepo.GetCommentByID(commentId)
	if err != nil {
		return Result{nil, err}
	}
	if userInfo.ID != comment.CreatedBy.ID.String() {
		return Result{nil, MakeUsecaseError("Isn't your comment!")}
	}

	err = u.placeRepo.DeleteComments(commentId)
	if err != nil {
		return Result{false, MakeUsecaseError("Cant delete comment(s).")}
	}
	return Result{true, nil}
}

func (u *PlaceUsecase) UpdateComment(id uuid.UUID, updateDto dto.UpdatePlaceCommentDto, userInfo dto.BaseUserInfo) Result {
	comment, err := u.placeRepo.GetCommentByID(id)
	if err != nil {
		return Result{nil, err}
	}
	if userInfo.ID != comment.CreatedBy.ID.String() {
		return Result{nil, MakeUsecaseError("Isn't your comment!")}
	}

	err = u.placeRepo.UpdateComment(id, updateDto)
	if err != nil {
		return Result{false, MakeUsecaseError("Cant update comment(s).")}
	}
	return Result{true, nil}
}
