package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/rohankarn35/nepsemarketbot/cmd"

	"github.com/rohankarn35/nepsemarketbot/server"
)

func main() {

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}
	// Connect to PostgreSQL
	db := cmd.InitializeDb(dsn)

	defer db.Close()

	ipoAPIURL := os.Getenv("IPO_API_URL")
	if ipoAPIURL == "" {
		log.Fatal("IPO_API_URL environment variable not set")
	}

	fpoAPIURL := os.Getenv("FPO_API_URL")
	if fpoAPIURL == "" {
		log.Fatal("FPO_API_URL environment variable not set")
	}

	c := cron.New(cron.WithLocation(time.FixedZone("NPT", 5*3600+45*60)))
	chatIDStr := os.Getenv("CHAT_ID")
	if chatIDStr == "" {
		log.Fatal("TELEGRAM_CHAT_ID environment variable not set")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Error converting TELEGRAM_CHAT_ID to int64: %v", err)
	}

	//initializebot
	bot := cmd.InitializeDataBase(botToken)

	server.InitializeScheduleronRestart(bot, c, db, chatID)

	// Add initial message to show bot is running
	log.Println("Bot started and waiting for messages...")
	err = cmd.SendMessages(db, c, ipoAPIURL, bot, fpoAPIURL, chatID)
	if err != nil {
		log.Printf("Error sending messages: %v", err)
	}

	// Periodic message sender

	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			err := cmd.SendMessages(db, c, ipoAPIURL, bot, fpoAPIURL, chatID)
			if err != nil {
				log.Printf("Error sending messages: %v", err)
				continue
			}

		}
	}()
	c.Start()

	// Keep the program running
	select {}
}
