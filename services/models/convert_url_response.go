package services_models

import "maps-to-waze-api/models"

type ConvertUrlResponse struct {
	URL         string             `json:"url"`
	Coordinates models.Coordinates `json:"coordinates"`
}
