package services_models

type GeoapifyReverseGeocodingResponse struct {
	Results []GRGResult `json:"results"`
}

type GRGResult struct {
	CountryCode  *string `json:"country_code"`
	Housenumber  *string `json:"housenumber"`
	Street       *string `json:"street"`
	Country      *string `json:"country"`
	Postcode     *string `json:"postcode"`
	State        *string `json:"state"`
	StateCode    *string `json:"state_code"`
	District     *string `json:"district"`
	City         *string `json:"city"`
	County       *string `json:"county"`
	CountyCode   *string `json:"county_code"`
	Formatted    *string `json:"formatted"`
	AddressLine1 *string `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
}
