package services

func ConvertDate(adDate, bsDate string) string {
	// Parse AD date into time.Time
	parsedADDate, _ := ParseEnglishMonth(adDate)

	parsedBSDate, _ := ParseNepaliDate(bsDate)

	// Format AD date to "JAN 20"

	// Combine with BS date
	return parsedBSDate + " ( " + parsedADDate + " )"
}
