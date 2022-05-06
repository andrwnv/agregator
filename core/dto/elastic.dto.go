package dto

const (
	EVENT_LOCATION_TYPE = "event"
	PLACE_LOCATION_TYPE = "place"
)

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

type CreateAggregatorRecordDto struct {
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
