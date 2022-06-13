package usecases

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
)

type SearchUsecase struct {
	placeRepo *repo.PlaceRepo
	eventRepo *repo.EventRepo
	esService *services.EsService
}

func NewSearchUsecase(
	placeRepo *repo.PlaceRepo,
	eventRepo *repo.EventRepo,
	esService *services.EsService) *SearchUsecase {
	return &SearchUsecase{
		placeRepo: placeRepo,
		eventRepo: eventRepo,
		esService: esService,
	}
}

func (u *SearchUsecase) Search(dto_ dto.SearchDto) Result {
	var result []dto.AggregatorRecordDto

	for _, objType := range dto_.SearchType {
		searchResult, _ := u.esService.Search(dto_.ValueToSearch, objType, int(dto_.From), int(dto_.Limit))
		result = append(result, searchResult...)
	}

	return Result{result, nil}
}

func (u *SearchUsecase) SearchNearby(dto_ dto.SearchNearbyDto) Result {
	nearby, err := u.esService.SearchNearby(dto_.Coords, int(dto_.From), int(dto_.Limit))

	var result []dto.AggregatorRecordDto
	for _, item := range nearby {
		for _, searchType := range dto_.SearchType {
			if searchType == item.LocationType {
				result = append(result, item)
			}
		}
	}

	return Result{result, err}
}
