package models

import "fmt"

type EventDuration struct {
	Minutes        int    `json:"minutes"`
	GenitiveMinute string `json:"genitive_minutes"`
}

type Event struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	MonthYear     string        `json:"month_year"`
	EventDate     string        `json:"event_date"`
	EventDuration EventDuration `json:"event_duration"`
	Image         string        `json:"image"`
	Route         string        `json:"route"`
}

func (event Event) String() string {
	return fmt.Sprintf("🎭*Нова вистава*: %s\n*Дата*: %s\n*Тривалість*: %d хв\n\n🔗*Посилання*: %s",
		event.Name, event.EventDate, event.EventDuration.Minutes, event.Route)
}

func (event Event) BackfillString() string {
	return fmt.Sprintf("🎭*Нова вистава*: %s\n*Дата*: %s\n\n🔗*Посилання*: %s",
		event.Name, event.EventDate, event.Route)
}
