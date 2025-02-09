package cmd

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/machinebox/graphql"
	"github.com/robfig/cron/v3"
	"github.com/rohankarn35/htmlcapture"
	ipodb "github.com/rohankarn35/nepsemarketbot/db"
	dbgraphql "github.com/rohankarn35/nepsemarketbot/graphql"
	"github.com/rohankarn35/nepsemarketbot/models"
	"github.com/rohankarn35/nepsemarketbot/server"
	"github.com/rohankarn35/nepsemarketbot/services"
)

func SendMessages(db *pgxpool.Pool, c *cron.Cron, bot *tgbotapi.BotAPI, chatID int64, client *graphql.Client) error {

	// Fetch latest data
	ipoData, fpoData, err := dbgraphql.GetIPOFPODetails(client)
	if err != nil {
		return err
	}
	log.Printf("Fetched IPO and FPO data successfully")

	var IPOdata []models.IPOAlertModel
	for _, ipo := range ipoData {
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

	for _, fpo := range fpoData {
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
			log.Print(isAvailable)
			if isAvailable {

				ipodb.CreateOrUpdateDB(db, ipo.CompanyName, ipo.StockSymbol, ipo.ShareType, ipo.SectorName, ipo.Status, ipo.PricePerUnit, ipo.MinUnits, ipo.MaxUnits, ipo.OpeningDateAD, ipo.OpeningDateBS, ipo.ClosingDateAD, ipo.ClosingDateBS, ipo.ClosingDateClosingTime, ipo.ShareRegistrar, ipo.Rating, ipo.Type)
				server.Scheduler(ipo.ClosingDateAD, ipo.ClosingDateClosingTime, ipo, bot, c, chatID)
				ipoCronData := models.CronJobModel{
					Closingdate: ipo.ClosingDateAD,
					Closingtime: ipo.ClosingDateClosingTime,
					Symbol:      ipo.StockSymbol,
				}
				ipodb.StoreCron(db, ipoCronData)
				ipoType := ipo.ShareType
				if ipo.ShareType == "ordinary" {
					ipoType = "General Public"
				}
				status := "Upcoming"
				if ipo.Status != "Nearing" {
					status = ipo.Status
				}
				openingDate := services.ConvertDate(ipo.OpeningDateAD, ipo.OpeningDateBS)
				closingDate := services.ConvertDate(ipo.ClosingDateAD, ipo.ClosingDateBS)
				opts := htmlcapture.CaptureOptions{
					Input: "templates/ipoAlert.html",
					Variables: map[string]string{
						"CompanyName": ipo.CompanyName,
						"Title":       status + " " + ipo.Type + " Alert",
						"Subtitle":    "(" + "For " + ipoType + ")",

						"IssueDate":   openingDate,
						"ClosingDate": closingDate,
						"IssuePrice":  "Rs. " + ipo.PricePerUnit,
						"Sector":      ipo.SectorName,
					},
					Selector:  ".container",
					ViewportW: 700,
					ViewportH: 600,
				}
				img, err := htmlcapture.Capture(opts)
				if err != nil {
					log.Fatalf("Error capturing screenshot: %v", err)
				}

				// Prepare the message text
				responseText := services.FormatIPOMessage(ipo)

				// Send the photo
				photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{Name: "ipoimage", Bytes: img})
				photo.Caption = responseText
				photo.ParseMode = "Markdown"

				// If IPO is open, add a button
				if strings.ToLower(ipo.Status) == "open" {
					button1 := tgbotapi.NewInlineKeyboardButtonURL("APPLY HERE", "https://meroshare.cdsc.com.np/")
					inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(button1),
					)
					photo.ReplyMarkup = inlineKeyboard
				}

				// Send the photo with caption and button
				if _, err := bot.Send(photo); err != nil {

					log.Printf("Error sending IPO image: %v", err)
					continue
				}

				log.Printf("Successfully sent IPO image to chat ID: %d", chatID)
			} else {
				log.Printf("Message Already Sent: %d", chatID)
			}
		} else {
			ipodb.DeleteIPOs(db, ipo.StockSymbol)
		}
	}
	return nil
}

func ScheduleSendMessage(db *pgxpool.Pool, c *cron.Cron, bot *tgbotapi.BotAPI, chatID int64, client *graphql.Client) {

	if err := SendMessages(db, c, bot, chatID, client); err != nil {
		log.Printf("Error sending messages: %v", err)
		if err := SendMessages(db, c, bot, chatID, client); err != nil {
			log.Printf("Error sending messages on retry: %v", err)
		}
	}

	_, err := c.AddFunc("0 9-17 * * *", func() {
		if err := SendMessages(db, c, bot, chatID, client); err != nil {
			log.Printf("Error sending messages: %v", err)
			if err := SendMessages(db, c, bot, chatID, client); err != nil {
				log.Printf("Error sending messages on retry: %v", err)
			}
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling send messages: %v", err)
	}
}
