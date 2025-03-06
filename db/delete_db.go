package ipodb

import (
	"fmt"
	"log"

	"github.com/rohankarn35/nepsemarketbot/db/models"
	"gorm.io/gorm"
)

func DeleteIPOs(db *gorm.DB, uniqueSymbol string) error {
	var count int64
	if err := db.Model(&models.NepseData{}).Where("unique_symbol=?", uniqueSymbol).Count(&count).Error; err != nil {
		return fmt.Errorf("failed checking symbol %v", err)

	}

	if count > 0 {
		result := db.Where("unique_symbol = ?", uniqueSymbol).
			Delete(&models.NepseData{})
		if result.Error != nil {
			log.Fatalf("Error deleting IPO: %v\n", result.Error)
			return result.Error
		}
		fmt.Println("IPO deleted successfully!")
	}
	return nil

}
