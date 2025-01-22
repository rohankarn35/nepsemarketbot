package ipodb

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rohankarn35/nepsemarketbot/models"
)

func ReadDB(db *pgxpool.Pool) {
	query := `SELECT * FROM nepsedata`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Error reading IPOs: %v\n", err)
	}
	defer rows.Close()

	fmt.Println("IPOs:")
	for rows.Next() {
		var ipo models.IPODB
		err := rows.Scan(
			&ipo.CompanyName,
			&ipo.StockSymbol,
			&ipo.ShareType,
			&ipo.SectorName,
			&ipo.Status,
			&ipo.PricePerUnit,
			&ipo.MinUnits,
			&ipo.MaxUnits,
			&ipo.OpeningDateAD,
			&ipo.OpeningDateBS,
			&ipo.ClosingDateAD,
			&ipo.ClosingDateBS,
			&ipo.ClosingDateClosingTime,
			&ipo.ShareRegistrar,
			&ipo.Rating,
			&ipo.Type,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("%+v\n", ipo)
	}
}
