package scraper

import (
	lastEvent "ticket-scrapper/internal/last-event"
	"ticket-scrapper/internal/models"
	"time"
)

func (c *Crawler) backfillEvent(id int) {
	timer := time.NewTimer(5 * time.Hour)
	defer timer.Stop()

	event, err := c.getEventByID(id)
	for err != nil {
		time.Sleep(15 * time.Second)

		select {
		case <-timer.C:
			return
		default:
			{
				event, err = c.getEventByID(id)
			}
		}
	}

	go c.bot.SendMessage(event.BackfillString())
}

func (c *Crawler) updateLastEventAndBackfill(events []models.Event) {
	for _, e := range events {
		currentLast := lastEvent.LastEvent.Load()
		if int64(e.ID) > currentLast {
			if e.ID-int(currentLast) > 1 {
				for id := int(currentLast) + 1; id < e.ID; id++ {
					go c.backfillEvent(id)
				}
			}
			lastEvent.LastEvent.Store(int64(e.ID))
		}
	}
}
