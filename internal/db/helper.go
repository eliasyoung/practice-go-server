package db

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseToPgTimestamptz(input string) (pgtype.Timestamptz, error) {
	// Parse string to Go's time.Time
	parsedTime, err := time.Parse(time.DateTime, input)
	log.Println(parsedTime)
	if err != nil {
		return pgtype.Timestamptz{}, err
	}

	// Convert time.Time to pgtype.Timestamptz
	timestamptz := pgtype.Timestamptz{
		Time:  parsedTime,
		Valid: true,
	}

	return timestamptz, nil
}
