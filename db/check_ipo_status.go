package ipodb

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckAndUpdateIPOStatus(db *pgxpool.Pool, uniqueSymbol string, status string) bool {
	var existingStatus string
	query := `SELECT IPOStatus FROM NepseData WHERE UniqueSymbol = $1`
	err := db.QueryRow(context.Background(), query, uniqueSymbol).Scan(&existingStatus)

	if err != nil {
		// If no rows are returned, insert the new IPO
		if err.Error() == "no rows in result set" {

			return true
		}
		log.Fatalf("Error querying IPO status: %v\n", err)
	}
	log.Printf("Existing Status: %s, New Status: %s, IPO Name: %s\n", existingStatus, status, uniqueSymbol)

	return existingStatus != status
}
