package models

//CalendarEvent ...
type CalendarEvent struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
