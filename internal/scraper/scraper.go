package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"ticket-scrapper/internal/bot"
	lastEvent "ticket-scrapper/internal/last-event"
	"ticket-scrapper/internal/models"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	url    string
	period time.Duration
	bot    *bot.Bot
}

func New(url string, period time.Duration, notificationBot *bot.Bot) *Crawler {
	return &Crawler{
		url:    url,
		period: period,
		bot:    notificationBot,
	}
}

func (c *Crawler) StartCrawler() {
	ticker := time.NewTicker(c.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			{
				fmt.Println("Parsing ", c.url)
				fmt.Println("Last event ID: ", lastEvent.LastEvent.Load())
				fmt.Println("Gorutines: ", runtime.NumGoroutine())
				go c.CrawlNewEventsAndNotify()
			}
		}
	}
}

func (c Crawler) CrawlNewEventsAndNotify() {
	events, err := c.CrawlNewEvents()
	if err != nil {
		log.Println("Error crawling events:", err)
		return
	}

	go c.updateLastEventAndBackfill(events)

	if len(events) > 0 {
		for _, event := range events {
			go c.bot.SendMessage(event.String())
			// fmt.Println(event.String())
		}
	}
}

// events sorted by id, it is crucial for backfill
func (c *Crawler) CrawlNewEvents() ([]models.Event, error) {
	totalEvents, err := c.getTotalEvents()
	if err != nil {
		return nil, err
	}

	events, err := c.getAllEvents(totalEvents)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(events, func(a models.Event, b models.Event) int {
		return a.ID - b.ID
	})

	return events, nil
}

func (c *Crawler) getTotalEvents() (int, error) {

	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 0, err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return 0, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return 0, err
	}

	var result models.ResponsePages
	err = json.Unmarshal([]byte(responseBody), &result)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return 0, err
	}

	return (result.Pagination.Total * result.Pagination.PerPage), nil
}
func (c *Crawler) getEventByID(ID int) (*models.Event, error) {
	url := fmt.Sprintf("%s/%d", c.url, ID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	eventInfo := doc.Find("div.ticketSelectionHead__left")
	if eventInfo.Length() == 0 {
		return nil, fmt.Errorf("event info not found")
	}

	event := &models.Event{
		ID: ID,
	}

	event.Name = eventInfo.Find(".ticketSelectionHead__name").Text()
	event.Route = url
	eventDate := eventInfo.Find(".ticketSelectionHead__item time").Text()
	if eventDate != "" {
		event.EventDate = eventDate
	}
	durationText := doc.Find(".performanceCard__time-val").Text()
	if durationText != "" {
		minutes, err := strconv.Atoi(strings.TrimSpace(durationText))
		if err == nil {
			event.EventDuration = models.EventDuration{Minutes: minutes}
		}
	}
	imgSrc, exists := doc.Find(".performanceCard__pic img").Attr("src")
	if exists {
		event.Image = imgSrc
	}

	return event, nil
}

func (c *Crawler) getAllEvents(totalEvents int) ([]models.Event, error) {
	url := fmt.Sprintf("%s?per_page=%d", c.url, totalEvents)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	var result models.ResponseEvents
	err = json.Unmarshal([]byte(responseBody), &result)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	events := filterSlice(result.Events, func(e models.Event) bool {
		currentLast := lastEvent.LastEvent.Load()

		return int64(e.ID) > currentLast
	})

	return events, nil
}

func filterSlice(events []models.Event, filter func(e models.Event) bool) []models.Event {
	var res []models.Event
	for _, event := range events {
		if filter(event) {
			res = append(res, event)
		}
	}

	return res
}
