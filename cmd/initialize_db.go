package cmd

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeDb(dbUrl string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	log.Print("Connected to Postgres")
	return db
}
