package models

//CalendarEvent represents an event on the calendar
type CalendarEvent struct {
	Title     string `json:"title"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
