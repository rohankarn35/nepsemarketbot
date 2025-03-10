package server

// import (
// 	"fmt"
// 	"log"
// 	"strings"
// 	"time"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/robfig/cron/v3"
// 	"github.com/rohankarn35/htmlcapture"
// 	"github.com/rohankarn35/nepsemarketbot/models"
// 	"github.com/rohankarn35/nepsemarketbot/services"
// )

// func Scheduler(closingdate, closingtime string, ipoData models.IPOAlertModel, bot *tgbotapi.BotAPI, c *cron.Cron, chatID int64) {
// 	// Parse the closing date and time
// 	layout := "2006-01-02 3:04 PM"
// 	closingDateTimeStr := fmt.Sprintf("%s %s", closingdate, closingtime)
// 	closingDateTime, err := time.Parse(layout, closingDateTimeStr)
// 	if err != nil {
// 		fmt.Println("Error parsing date and time:", err)
// 		return
// 	}

// 	// Subtract one hour from the closing time
// 	reminderTime := closingDateTime.Add(-1 * time.Hour)

// 	// Format the reminder time for cron
// 	cronSchedule := fmt.Sprintf("%d %d %d %d *", reminderTime.Minute(), reminderTime.Hour(), reminderTime.Day(), reminderTime.Month())

// 	// Nepal Time Zone (UTC+5:45)

// 	_, err = c.AddFunc(cronSchedule, func() {
// 		RemainderFunction(ipoData, bot, chatID)
// 	})
// 	if err != nil {
// 		fmt.Println("Error scheduling task:", err)
// 		return
// 	}
// 	fmt.Println("The scheduler has beeen done")
// 	// Start the cron scheduler

// }

// func RemainderFunction(ipo models.IPOAlertModel, bot *tgbotapi.BotAPI, chatID int64) {
// 	ipoType := ipo.ShareType
// 	if ipo.ShareType == "ordinary" {
// 		ipoType = "General Public"
// 	}
// 	status := "Upcoming"
// 	if ipo.Status != "Nearing" {
// 		status = ipo.Status
// 	}
// 	openingDate := services.ConvertDate(ipo.OpeningDateAD, ipo.OpeningDateBS)
// 	closingDate := services.ConvertDate(ipo.ClosingDateAD, ipo.ClosingDateBS)
// 	opts := htmlcapture.CaptureOptions{
// 		Input: "templates/ipowarning.html",
// 		Variables: map[string]string{
// 			"CompanyName": ipo.CompanyName,
// 			"Title":       status + " " + ipo.Type + " Alert",
// 			"Subtitle":    "(" + "For " + ipoType + ")",

// 			"IssueDate":   openingDate,
// 			"ClosingDate": closingDate,
// 			"IssuePrice":  "Rs. " + ipo.PricePerUnit,
// 			"Sector":      ipo.SectorName,
// 		},
// 		Selector:  ".container",
// 		ViewportW: 600,
// 		ViewportH: 600,
// 	}
// 	img, err := htmlcapture.Capture(opts)
// 	if err != nil {
// 		log.Fatalf("Error capturing screenshot: %v", err)
// 	}

// 	// Prepare the message text
// 	responseText := services.FormatIPOAlertMessage(ipo)

// 	// Send the photo
// 	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{Name: "ipoimage", Bytes: img})
// 	photo.Caption = responseText
// 	photo.ParseMode = "Markdown"

// 	// If IPO is open, add a button
// 	if strings.ToLower(ipo.Status) == "open" {
// 		button1 := tgbotapi.NewInlineKeyboardButtonURL("APPLY HERE", "https://meroshare.cdsc.com.np/")
// 		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
// 			tgbotapi.NewInlineKeyboardRow(button1),
// 		)
// 		photo.ReplyMarkup = inlineKeyboard
// 	}

// 	log.Printf("Attempting to send IPO message to chat ID: %d", chatID)

// 	// Try to send the message up to 2 times if an error occurs
// 	for i := 0; i < 2; i++ {
// 		if _, err := bot.Send(photo); err != nil {
// 			log.Printf("Error sending IPO message (attempt %d): %v", i+1, err)
// 			if i == 1 {
// 				log.Printf("Failed to send IPO message after 2 attempts")
// 				return
// 			}
// 		} else {
// 			log.Printf("Successfully sent IPO message to chat ID: %d", chatID)
// 			break
// 		}
// 	}

// }
// func InitializeScheduleronRestart(bot *tgbotapi.BotAPI, c *cron.Cron, db *pgxpool.Pool, chatID int64) {
// 	// ipoCronData, err := ipodb.ReadCron(db)
// 	// if err != nil {
// 	// 	log.Printf("Error reading cron data from database: %v", err)
// 	// 	return
// 	// }

// 	// for _, cronJob := range ipoCronData {
// 	// 	Scheduler(cronJob.Closingdate, cronJob.Closingtime, cronJob.IPOAlertModel, bot, c, chatID)
// 	// }

// }
