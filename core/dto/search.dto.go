package dto

type SearchDto struct {
	ValueToSearch string   `json:"value"`
	SearchType    []string `json:"search_types"` // event, place
	From          uint     `json:"from"`
	Limit         uint     `json:"limit"`
}

type SearchNearbyDto struct {
	Coords     LocationDto `json:"coords"`
	SearchType []string    `json:"search_types"` // event, place
	From       uint        `json:"from"`
	Limit      uint        `json:"limit"`
}
