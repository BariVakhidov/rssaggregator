package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/BariVakhidov/rssaggregator/internal/database"
)

func ConnectToDB() (*database.Queries, error) {
	dbPortString := os.Getenv("DB_URL")

	if dbPortString == "" {
		return nil, fmt.Errorf("no DB_URL found")
	}

	dbConnection, dbErr := sql.Open("postgres", dbPortString)
	if dbErr != nil {
		return nil, fmt.Errorf("can't connect to db: %v", dbErr)
	}

	db := database.New(dbConnection)

	return db, nil
}
