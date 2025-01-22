package ipodb

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateOrUpdateDB(db *pgxpool.Pool, companyName, stockSymbol, shareType, sectorName, ipoStatus, pricePerUnit, minUnits, maxUnits, openingDateAD, openingDateBS, closingDateAD, closingDateBS, closingDateClosingTime, shareRegistrar, rating, types string) {
	query := `
		INSERT INTO nepsedata (CompanyName, StockSymbol, ShareType, SectorName, IPOStatus, PricePerUnit, MinUnits, MaxUnits, OpeningDateAD, OpeningDateBS, ClosingDateAD, ClosingDateBS, ClosingDateClosingTime, ShareRegistrar, Rating, Types)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (StockSymbol) DO UPDATE SET
			CompanyName = EXCLUDED.CompanyName,
			ShareType = EXCLUDED.ShareType,
			SectorName = EXCLUDED.SectorName,
			IPOStatus = EXCLUDED.IPOStatus,
			PricePerUnit = EXCLUDED.PricePerUnit,
			MinUnits = EXCLUDED.MinUnits,
			MaxUnits = EXCLUDED.MaxUnits,
			OpeningDateAD = EXCLUDED.OpeningDateAD,
			OpeningDateBS = EXCLUDED.OpeningDateBS,
			ClosingDateAD = EXCLUDED.ClosingDateAD,
			ClosingDateBS = EXCLUDED.ClosingDateBS,
			ClosingDateClosingTime = EXCLUDED.ClosingDateClosingTime,
			ShareRegistrar = EXCLUDED.ShareRegistrar,
			Rating = EXCLUDED.Rating,
			Types = EXCLUDED.Types
		WHERE nepsedata.IPOStatus IS DISTINCT FROM EXCLUDED.IPOStatus
	`
	_, err := db.Exec(context.Background(), query, companyName, stockSymbol, shareType, sectorName, ipoStatus, pricePerUnit, minUnits, maxUnits, openingDateAD, openingDateBS, closingDateAD, closingDateBS, closingDateClosingTime, shareRegistrar, rating, types)
	if err != nil {
		log.Fatalf("Error inserting or updating IPO: %v\n", err)
	}
	fmt.Println("IPO created or updated successfully!")
}
