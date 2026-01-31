package persistance

import (
	"database/sql"
	"fmt"
	"task-management/user-service/src/internal/config"
	errors "task-management/user-service/src/pkg/error"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func ConnectToDatabase(config *config.Config) (*Database, error) {

	databaseUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		fmt.Printf("Error is :- %v", err)
		return nil, errors.ErrInternalServer
	}

	fmt.Printf("Connected to database \n")

	return &Database{db: db}, nil

}
