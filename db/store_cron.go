package ipodb

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rohankarn35/nepsemarketbot/models"
)

func StoreCron(db *pgxpool.Pool, cron models.CronJobModel) error {
	query := `
		INSERT INTO cron_jobs (Closingdate, Closingtime, StockSymbol)
		VALUES ($1, $2, $3)
		ON CONFLICT (StockSymbol) DO NOTHING
		RETURNING JobID;
	`
	var jobID int
	ctx := context.Background()
	err := db.QueryRow(ctx, query, cron.Closingdate, cron.Closingtime, cron.Symbol).Scan(&jobID)
	if err != nil {
		return fmt.Errorf("failed to insert cron job: %w", err)
	}

	log.Printf("Cron job inserted successfully with JobID: %d", jobID)
	return nil

}

func ReadCron(db *pgxpool.Pool) ([]models.CronJobIpoModel, error) {
	query := `
		SELECT cj.Closingdate, cj.Closingtime, cj.StockSymbol, nd.*
		FROM cron_jobs cj
		JOIN nepsedata nd ON cj.StockSymbol = nd.StockSymbol;
	`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to read cron jobs: %w", err)
	}
	defer rows.Close()

	var cronJobs []models.CronJobIpoModel
	for rows.Next() {
		var cronJob models.CronJobIpoModel
		err := rows.Scan(
			&cronJob.Closingdate,
			&cronJob.Closingtime,
			&cronJob.Symbol,
			&cronJob.CompanyName,
			&cronJob.StockSymbol,
			&cronJob.ShareType,
			&cronJob.SectorName,
			&cronJob.Status,
			&cronJob.PricePerUnit,
			&cronJob.MinUnits,
			&cronJob.MaxUnits,
			&cronJob.OpeningDateAD,
			&cronJob.OpeningDateBS,
			&cronJob.ClosingDateAD,
			&cronJob.ClosingDateBS,
			&cronJob.ClosingDateClosingTime,
			&cronJob.ShareRegistrar,
			&cronJob.Rating,
			&cronJob.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cron job: %w", err)
		}
		cronJobs = append(cronJobs, cronJob)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return cronJobs, nil
}
