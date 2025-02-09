package ipodb

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckAndUpdateIPOStatus(db *pgxpool.Pool, symbol string, status string) bool {
	var existingStatus string
	query := `SELECT IPOStatus FROM nepsedata WHERE StockSymbol = $1`
	err := db.QueryRow(context.Background(), query, symbol).Scan(&existingStatus)

	if err != nil {
		// If no rows are returned, insert the new IPO
		if err.Error() == "no rows in result set" {

			return true
		}
		log.Fatalf("Error querying IPO status: %v\n", err)
	}
	log.Printf("Existing Status: %s, New Status: %s, IPO Name: %s\n", existingStatus, status, symbol)

	if existingStatus == status {
		return false
	}

	return true
}
