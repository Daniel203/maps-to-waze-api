package database

import (
	"context"
	"fmt"
	"time"
)

func GetNumberOfRequestsThisMonth(ctx context.Context) (int, error) {
	db, err := DBFromContext(ctx)

	if err != nil {
		return 0, err
	}

	monthStartDate := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	monthEndDate := monthStartDate.AddDate(0, 1, 0)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM request WHERE created_at >= $1 AND created_at < $2", monthStartDate, monthEndDate).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func InsertRequest(ctx context.Context, request_id string) error {
	if request_id == "" {
		return fmt.Errorf("request_id cannot be empty")
	}

	db, err := DBFromContext(ctx)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO request (http_request_id) VALUES ($1)", request_id)
	if err != nil {
		return err
	}

	return nil
}
