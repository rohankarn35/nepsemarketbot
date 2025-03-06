package ipodb

import (
	"fmt"
	"log"

	gorm_model "github.com/rohankarn35/nepsemarketbot/db/models"

	"gorm.io/gorm"
)

func StoreCron(db *gorm.DB, cron gorm_model.CronJob) error {

	if err := db.Create(&cron).Error; err != nil {
		return fmt.Errorf("failed to store cron job %v", err)
	}
	return nil

}

func ReadCron(db *gorm.DB) ([]gorm_model.CronJob, error) {
	var cron []gorm_model.CronJob

	if err := db.Preload("nepse_data").Find(&cron).Error; err != nil {
		return nil, fmt.Errorf("failed to load data")
	}
	log.Print("read all the documents", cron)
	return cron, nil
}
