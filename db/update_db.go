package ipodb

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateStatus(db *pgxpool.Pool, symbol string, status string) {
	query := `
		UPDATE nepseData 
		SET status = $1
		WHERE StockSymbol = $2
	`
	_, err := db.Exec(context.Background(), query, status, symbol)
	if err != nil {
		log.Fatalf("Error updating IPO status: %v\n", err)
	}
	fmt.Println("IPO status updated successfully!")
}
