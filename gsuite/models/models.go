package models

//CalendarEvent represents an event on the calendar
type CalendarEvent struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
