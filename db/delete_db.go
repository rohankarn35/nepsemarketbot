package ipodb

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DeleteIPOs(db *pgxpool.Pool, symbol string) {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM nepsedata WHERE StockSymbol = $1)`
	err := db.QueryRow(context.Background(), checkQuery, symbol).Scan(&exists)
	if err != nil {
		log.Fatalf("Error checking if symbol exists: %v\n", err)
	}

	if exists {
		query := `DELETE FROM nepseData WHERE StockSymbol = $1`
		_, err := db.Exec(context.Background(), query, symbol)
		if err != nil {
			log.Fatalf("Error deleting IPO: %v\n", err)
		}
		fmt.Println("IPO deleted successfully!")
	}
}
