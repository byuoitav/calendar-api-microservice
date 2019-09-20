package models

//GSuiteEvent stores an event in G Suite format
type GSuiteEvent struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

//CalendarEvent stores an event
type CalendarEvent struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

//CouchScheduleConfig stores the scheduling panel configuration
type CouchScheduleConfig struct {
	ID           string `json:"_id"`
	Rev          string `json:"_rev"`
	Resource     string `json:"resource"`
	Name         string `json:"displayname"`
	AutoDiscURL  string `json:"autodiscover-url"`
	AccessType   string `json:"access-type"`
	Image        string `json:"image-url"`
	BookNow      bool   `json:"allowbooknow"`
	ShowHelp     bool   `json:"showhelp"`
	CalendarType string `json:"calendar-type"`
}
