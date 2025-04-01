package services

import (
	"context"
	"fmt"
	"maps-to-waze-api/internal/database"
)

func checkNumberOfRequestsThisMonth(ctx context.Context, requestTypeId int, multiplier *float64, requestLimit int) (bool, error) {
	if multiplier == nil {
		multiplier = new(float64)
		*multiplier = 1.0
	}

	requests, err := database.GetNumberOfRequestsThisMonth(ctx, requestTypeId)
	if err != nil {
		return false, fmt.Errorf("failed to get the number of requests this month: %w", err)
	}

	return float64(requests) * (*multiplier) < float64(requestLimit), nil
}

func checkNumberOfRequestsToday(ctx context.Context, requestTypeId int, multiplier *float64, requestLimit int) (bool, error) {
	if multiplier == nil {
		multiplier = new(float64)
		*multiplier = 1.0
	}
	
	requests, err := database.GetNumberOfRequestsToday(ctx, requestTypeId)
	if err != nil {
		return false, fmt.Errorf("failed to get the number of requests today: %w", err)
	}

	return float64(requests) * (*multiplier) < float64(requestLimit), nil
}
