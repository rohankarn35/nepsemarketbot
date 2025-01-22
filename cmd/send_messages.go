package cmd

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	ipodb "github.com/rohankarn35/nepsemarketbot/db"
	"github.com/rohankarn35/nepsemarketbot/models"
	"github.com/rohankarn35/nepsemarketbot/server"
	"github.com/rohankarn35/nepsemarketbot/services"
)

func SendMessages(db *pgxpool.Pool, c *cron.Cron, ipoAPIURL string, bot *tgbotapi.BotAPI, fpoAPIURL string, chatID int64) error {

	// Fetch latest data
	ipoData, err := server.FetchIPOData(ipoAPIURL)
	if err != nil {
		return err

	}
	log.Printf("Fetched IPO data successfully")

	fpoData, err := server.FetchFPOData(fpoAPIURL)
	if err != nil {
		return err

	}
	log.Printf("Fetched FPO data successfully")
	var IPOdata []models.IPOAlertModel
	for _, ipo := range ipoData.Result.Data {
		IPOdata = append(IPOdata, models.IPOAlertModel{
			CompanyName:            ipo.CompanyName,
			StockSymbol:            ipo.StockSymbol,
			ShareType:              ipo.ShareType,
			Status:                 ipo.Status,
			PricePerUnit:           ipo.PricePerUnit,
			MinUnits:               ipo.MinUnits,
			MaxUnits:               ipo.MaxUnits,
			OpeningDateAD:          ipo.OpeningDateAD,
			OpeningDateBS:          ipo.OpeningDateBS,
			ClosingDateAD:          ipo.ClosingDateAD,
			ClosingDateBS:          ipo.ClosingDateBS,
			ClosingDateClosingTime: ipo.ClosingDateClosingTime,
			ShareRegistrar:         ipo.ShareRegistrar,
			Rating:                 ipo.Rating,
			SectorName:             ipo.SectorName,
			Type:                   "IPO",
		})
	}

	for _, fpo := range fpoData.Result.Data {
		IPOdata = append(IPOdata, models.IPOAlertModel{
			CompanyName:            fpo.CompanyName,
			StockSymbol:            fpo.StockSymbol,
			ShareType:              fpo.ShareType,
			Status:                 fpo.Status,
			PricePerUnit:           fpo.PricePerUnit,
			MinUnits:               fpo.MinUnits,
			MaxUnits:               fpo.MaxUnits,
			OpeningDateAD:          fpo.OpeningDateAD,
			OpeningDateBS:          fpo.OpeningDateBS,
			ClosingDateAD:          fpo.ClosingDateAD,
			ClosingDateBS:          fpo.ClosingDateBS,
			ClosingDateClosingTime: fpo.ClosingDateClosingTime,
			ShareRegistrar:         fpo.ShareRegistrar,
			Rating:                 fpo.Rating,
			SectorName:             fpo.SectorName,
			Type:                   "FPO",
		})
	}

	// Send IPO updates
	for _, ipo := range IPOdata {
		if strings.ToLower(ipo.Status) == "nearing" || strings.ToLower(ipo.Status) == "open" {
			isAvailable := ipodb.CheckAndUpdateIPOStatus(db, ipo.StockSymbol, ipo.Status)
			if isAvailable {

				ipodb.CreateOrUpdateDB(db, ipo.CompanyName, ipo.StockSymbol, ipo.ShareType, ipo.SectorName, ipo.Status, ipo.PricePerUnit, ipo.MinUnits, ipo.MaxUnits, ipo.OpeningDateAD, ipo.OpeningDateBS, ipo.ClosingDateAD, ipo.ClosingDateBS, ipo.ClosingDateClosingTime, ipo.ShareRegistrar, ipo.Rating, ipo.Type)
				server.Scheduler(ipo.ClosingDateAD, ipo.ClosingDateClosingTime, ipo, bot, c, chatID)
				ipoCronData := models.CronJobModel{
					Closingdate: ipo.ClosingDateAD,
					Closingtime: ipo.ClosingDateClosingTime,
					Symbol:      ipo.StockSymbol,
				}
				ipodb.StoreCron(db, ipoCronData)

				responseText := services.FormatIPOMessage(ipo)
				msg := tgbotapi.NewMessage(chatID, responseText)
				if strings.ToLower(ipo.Status) == "open" {
					button1 := tgbotapi.NewInlineKeyboardButtonURL("APPLY HERE", "https://meroshare.cdsc.com.np/")
					inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(button1),
					)
					msg.ReplyMarkup = inlineKeyboard
				}

				msg.ParseMode = "Markdown"

				log.Printf("Attempting to send IPO message to chat ID: %d", chatID)

				// Actually send the message
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Error sending IPO message: %v", err)
					continue
				}
				log.Printf("Successfully sent IPO message to chat ID: %d", chatID)
			} else {
				log.Printf("Message Already Sent: %d", chatID)
			}
		} else {
			ipodb.DeleteIPOs(db, ipo.StockSymbol)
		}
	}
	return nil
}
