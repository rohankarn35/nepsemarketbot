package models

type Root struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Result     Result `json:"result"`
}

// Result represents the "result" field in the JSON
type Result struct {
	Data  []IPOData `json:"data"`
	Pager Pager     `json:"pager"`
}

// IPOData represents each IPO record in the "data" array
type IPOData struct {
	IpoID                   int    `json:"ipoId"`
	CompanyName             string `json:"companyName"`
	StockSymbol             string `json:"stockSymbol"`
	ShareRegistrar          string `json:"shareRegistrar"`
	SectorName              string `json:"sectorName"`
	FileName                string `json:"fileName"`
	ShareType               string `json:"shareType"`
	PricePerUnit            string `json:"pricePerUnit"`
	Rating                  string `json:"rating"`
	Units                   string `json:"units"`
	MinUnits                string `json:"minUnits"`
	MaxUnits                string `json:"maxUnits"`
	LocalUnits              string `json:"localUnits"`
	GeneralUnits            string `json:"generalUnits"`
	PromoterUnits           string `json:"promoterUnits"`
	MutualFundUnits         string `json:"mutualFundUnits"`
	OtherUnits              string `json:"otherUnits"`
	TotalAmount             string `json:"totalAmount"`
	OpeningDateAD           string `json:"openingDateAD"`
	OpeningDateBS           string `json:"openingDateBS"`
	ClosingDateAD           string `json:"closingDateAD"`
	ClosingDateBS           string `json:"closingDateBS"`
	ClosingDateClosingTime  string `json:"closingDateClosingTime"`
	ExtendedDateAD          string `json:"extendedDateAD"`
	ExtendedDateBS          string `json:"extendedDateBS"`
	ExtendedDateClosingTime string `json:"extendedDateClosingTime"`
	Status                  string `json:"status"`
	FiscalYearAD            string `json:"fiscalYearAD"`
	FiscalYearBS            string `json:"fiscalYearBS"`
	CultureCode             string `json:"cultureCode"`
}

// Pager represents the "pager" field in the JSON
type Pager struct {
	PageNo         int `json:"pageNo"`
	ItemsPerPage   int `json:"itemsPerPage"`
	PagePerDisplay int `json:"pagePerDisplay"`
	TotalNextPages int `json:"totalNextPages"`
}
