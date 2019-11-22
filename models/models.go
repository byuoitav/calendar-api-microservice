package models

//GSuiteEvent stores an event in G Suite format
type GSuiteEvent struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

//CalendarEvent stores an event
type CalendarEvent struct {
	Title     string `json:"title"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
