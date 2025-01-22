package models

type IPOAlertModel struct {
	CompanyName            string
	StockSymbol            string
	ShareType              string
	SectorName             string
	Status                 string
	PricePerUnit           string
	MinUnits               string
	MaxUnits               string
	OpeningDateAD          string
	OpeningDateBS          string
	ClosingDateAD          string
	ClosingDateBS          string
	ClosingDateClosingTime string
	ShareRegistrar         string
	Rating                 string
	Type                   string
}
