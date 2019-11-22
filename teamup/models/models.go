package models

//CalendarEvent represents an event on the calendar
type CalendarEvent struct {
	Title     string `json:"title"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

//TeamUpEventReturn represents the event data returned by the TeamUp API
type TeamUpEventReturn struct {
	ID                 string `json:"id"`
	RemoteID           string `json:"remote_id"`
	SeriesID           string `json:"series_id"`
	SubCalendarIDArray []int  `json:"subcalendar_id"`
	SubCalendarID      int    `json:"subcalendar_id"`
	StartDate          string `json:"start_dt"`
	EndDate            string `json:"end_dt"`
	AllDay             bool   `json:"all_day"`
	Title              string `json:"title"`
	Who                string `json:"who"`
	Location           string `json:"location"`
	Notes              string `json:"notes"`
	RecurRule          string `json:"rrule"`
	RecurInstStart     string `json:"ristart_dt"`
	RecurSeriesStart   string `json:"rsstart_dt"`
	TimeZone           string `json:"tz"`
	Version            string `json:"version"`
	ReadOnly           bool   `json:"readonly"`
	CreationDate       string `json:"creation_dt"`
	UpdateDate         string `json:"update_dt"`
}

//TeamUpEventResponse represents the json object returned by the TeamUp API events get request
type TeamUpEventResponse struct {
	Events    []TeamUpEventReturn `json:"events"`
	TimeStamp int                 `json:"timestamp"`
}

//TeamUpEventSend represents the event data to be sent to the TeamUp API
type TeamUpEventSend struct {
	SubCalendarID int    `json:"subcalendar_id"`
	StartDate     string `json:"start_dt"`
	EndDate       string `json:"end_dt"`
	AllDay        bool   `json:"all_day"`
	RecurRule     string `json:"rrule"`
	Title         string `json:"title"`
	Who           string `json:"who"`
	Location      string `json:"location"`
	Notes         string `json:"notes"`
}

//TeamUpSubcalendar represents a subcalendar in json
type TeamUpSubcalendar struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	Color        int    `json:"color"`
	Overlap      bool   `json:"overlap"`
	ReadOnly     bool   `json:"readonly"`
	CreationDate string `json:"creation_dt"`
	UpdateDate   string `json:"update_dt"`
}

//TeamUpSubcalendarList holds a list of subcalendars
type TeamUpSubcalendarList struct {
	Subcalendars []TeamUpSubcalendar `json:"subcalendars"`
}
