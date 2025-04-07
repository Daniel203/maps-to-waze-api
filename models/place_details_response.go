package models;

type PlaceDetailsResponse struct {
	Formatted *string `json:"formatted"`
	AddressLine1 *string `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
}

