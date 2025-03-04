package services_models

type GooglePlacesResponse struct {
	Status string `json:"status"`
	Result Result `json:"result"`
}

type Result struct {
	Geometry Geometry `json:"geometry"`
}

type Geometry struct {
	Location `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
