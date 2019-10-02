package models

// CalendarEvent models a calendar event
type CalendarEvent struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
