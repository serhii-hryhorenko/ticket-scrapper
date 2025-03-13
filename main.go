package main

import (
	"ticket-scrapper/internal/bot"
	lastEvent "ticket-scrapper/internal/last-event"
	"ticket-scrapper/internal/scraper"
	"time"
)

const URL = "https://sales.ft.org.ua/events"

// placeholder values for now
const botToken = ""
const channelID = 123

func main() {
	bot := bot.New(botToken, channelID)
	scraper := scraper.New(URL, 30*time.Second, &bot)

	lastEvent.LastEvent.Store(int64(5470))
	lastEvent.InitLastEvent()
	scraper.StartCrawler()
}
