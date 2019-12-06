package models

// CalendarEvent models a calendar event
type CalendarEvent struct {
	Title     string `json:"title"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// ExchangeToken models an exchange token
type ExchangeToken struct {
	Type       string `json:"token_type"`
	ExpireTime int    `json:"expires_in"`
	Token      string `json:"access_token"`
}

type ExchangeEventResponse struct {
	Events []ExchangeEvent `json:"value"`
}

// ExchangeEvent models an event returned by microsoft exchange service
type ExchangeEvent struct {
	ID                         string   `json:"id"`
	CreatedDateTime            string   `json:"createdDateTime"`
	LastModifiedDateTime       string   `json:"lastModifiedDateTime"`
	ChangeKey                  string   `json:"changeKey"`
	Categories                 []string `json:"categories"`
	OriginalStartTimeZone      string   `json:"originalStartTimeZone"`
	OriginalEndTimeZone        string   `json:"originalEndTimeZone"`
	ICalUID                    string   `json:"iCalUId"`
	ReminderMinutesBeforeStart int      `json:"reminderMinutesBeforeStart"`
	IsReminderOn               bool     `json:"isReminderOn"`
	HasAttachments             bool     `json:"hasAttachments"`
	Subject                    string   `json:"subject"`
	BodyPreview                string   `json:"bodyPreview"`
	Importance                 string   `json:"importance"`
	Sensitivity                string   `json:"sensitivity"`
	IsAllDay                   bool     `json:"isAllDay"`
	IsCancelled                bool     `json:"isCancelled"`
	IsOrganizer                bool     `json:"isOrganizer"`
	ResponseRequested          bool     `json:"responseRequested"`
	SeriesMasterID             string   `json:"seriesMasterId"`
	ShowAs                     string   `json:"showAs"`
	EventType                  string   `json:"type"`
	WebLink                    string   `json:"webLink"`
	OnlineMeetingURL           string   `json:"onlineMeetingUrl"`
	Recurrence                 string   `json:"recurrence"`
	// ResponseStatus             string       `json:"responseStatus"`
	Body  ExchangeBody `json:"body"`
	Start ExchangeDate `json:"start"`
	End   ExchangeDate `json:"end"`
	// Location  string       `json:"location"`
	// Locations string       `json:"locations"`
	// Attendees string       `json:"attendees"`
	// Organizer string       `json:"organizer"`
}

type ExchangeEventRequest struct {
	Subject   string             `json:"Subject"`
	Body      ExchangeBody       `json:"Body"`
	Start     ExchangeDate       `json:"Start"`
	End       ExchangeDate       `json:"End"`
	Attendees []ExchangeAttendee `json:"Attendees"`
}

type ExchangeBody struct {
	ContentType string `json:"ContentType"`
	Content     string `json:"Content"`
}

type ExchangeDate struct {
	DateTime string `json:"DateTime"`
	TimeZone string `json:"TimeZone"`
}

type ExchangeAttendee struct {
	EmailAddress ExchangeEmailAddress `json:"EmailAddress"`
	Type         string               `json:"Type"`
}

type ExchangeEmailAddress struct {
	Address string `json:"Address"`
	Name    string `json:"Name"`
}

type ExchangeCalenderResponse struct {
	Calendars []ExchangeCalendar `json:"value"`
}

type ExchangeCalendar struct {
	ID                            string                `json:"id"`
	Name                          string                `json:"name"`
	Color                         string                `json:"color"`
	IsDefault                     bool                  `json:"isDefaultCalendar"`
	ChangeKey                     string                `json:"changeKey"`
	CanShare                      bool                  `json:"canShare"`
	CanViewPrivate                bool                  `json:"canViewPrivateItems"`
	CanEdit                       bool                  `json:"canEdit"`
	AllowedOnlineMeetingProviders []string              `json:"allowedOnlineMeetingProviders"`
	DefaultOnlineMeetingProvider  string                `json:"defaultOnlineMeetingProvider"`
	TallyingResponses             bool                  `json:"isTallyingResponses"`
	Removable                     bool                  `json:"isRemovable"`
	Owner                         ExchangeCalendarOwner `json:"owner"`
}

type ExchangeCalendarOwner struct {
	Name    string `json:"name"`
	Address string `json:"Address"`
}
