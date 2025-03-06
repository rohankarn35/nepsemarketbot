package models

import "gorm.io/gorm"

type NepseData struct {
	gorm.Model                    // Adds ID, CreatedAt, UpdatedAt, DeletedAt fields
	UniqueSymbol           string `gorm:"unique;not null;primaryKey"` // Make it explicit primary key
	CompanyName            string
	StockSymbol            string
	ShareRegistrar         string
	SectorName             string
	ShareType              string
	PricePerUnit           string
	Rating                 string
	Units                  string
	MinUnits               string
	MaxUnits               string
	TotalAmount            string
	OpeningDateAD          string
	OpeningDateBS          string
	ClosingDateAD          string
	ClosingDateBS          string
	ClosingDateClosingTime string
	Status                 string
}

type CronJob struct {
	gorm.Model                    // Adds ID, CreatedAt, UpdatedAt, DeletedAt fields
	UniqueSymbol           string `gorm:"unique;not null"` // Foreign key, not primary key
	StockSymbol            string
	OpeningDateAD          string
	OpeningDateBS          string
	ClosingDateAD          string
	ClosingDateBS          string
	ClosingDateClosingTime string
	Status                 string
	NepseData              NepseData `gorm:"foreignKey:UniqueSymbol;references:UniqueSymbol"` // One-to-one relationship
}
