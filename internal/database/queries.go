package database

import (
	"context"
	"fmt"
	"time"
)

func GetNumberOfRequestsThisMonth(ctx context.Context, requestTypeId int) (int, error) {
	db, err := DBFromContext(ctx)

	if err != nil {
		return 0, err
	}

	monthStartDate := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	monthEndDate := monthStartDate.AddDate(0, 1, 0)

	var count int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM request WHERE request_type_id = $1 AND created_at >= $2 AND created_at < $3",
		requestTypeId,
		monthStartDate,
		monthEndDate,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetNumberOfRequestsToday(ctx context.Context, requestTypeId int) (int, error) {
	db, err := DBFromContext(ctx)

	if err != nil {
		return 0, err
	}

	todayStartDate := time.Now().Truncate(24 * time.Hour)
	todayEndDate := todayStartDate.AddDate(0, 0, 1)

	var count int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM request WHERE request_type_id = $1 AND created_at >= $2 AND created_at < $3",
		requestTypeId,
		todayStartDate,
		todayEndDate,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func InsertRequest(ctx context.Context, requestId string, requestTypeId int) error {
	if requestId == "" {
		return fmt.Errorf("request_id cannot be empty")
	}

	db, err := DBFromContext(ctx)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO request (http_request_id, request_type_id) VALUES ($1, $2)",
		requestId,
		requestTypeId,
	)

	return err
}
