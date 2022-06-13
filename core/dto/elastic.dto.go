package dto

import "github.com/google/uuid"

type LocationDto struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type AggregatorRecordDto struct {
	ID           string      `json:"id"`
	LocationName string      `json:"location_name"`
	Location     LocationDto `json:"location"`
	LocationType string      `json:"location_type"`
}

type AggregatorRecordElasticDto struct {
	ID                  string              `json:"_id"`
	AggregatorRecordDto AggregatorRecordDto `json:"_source"`
}

type CreateAggregatorRecordDto struct {
	ID           uuid.UUID   `json:"id"`
	LocationName string      `json:"location_name"`
	Location     LocationDto `json:"location"`
	LocationType string      `json:"location_type"`
}

type UpdateAggregatorRecordDto struct {
	LocationName string      `json:"location_name"`
	Location     LocationDto `json:"location"`
	LocationType string      `json:"location_type"`
}

type GeoBounding struct {
	TopLeft     LocationDto `json:"top_left"`
	BottomRight LocationDto `json:"bottom_right"`
}

type GeoSearch struct {
	GeoBoundingBox GeoBounding `json:"geo_bounding_box"`
}
