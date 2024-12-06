package db

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func SetupLogger(l *zap.SugaredLogger) {
	logger = l
}

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
