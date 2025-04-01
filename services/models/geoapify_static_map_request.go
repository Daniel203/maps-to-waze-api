package services_models

type GeoapifyStaticMapRequest struct {
	Style       string   `json:"style"`
	ScaleFactor int      `json:"scaleFactor"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Center      Center   `json:"center"`
	Zoom        int      `json:"zoom"`
	Markers     []Marker `json:"markers"`
}

type Marker struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Color string  `json:"color"`
	Size  string  `json:"size"`
}

type Center struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
