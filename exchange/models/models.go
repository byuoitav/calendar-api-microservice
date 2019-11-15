package models

// CalendarEvent models a calendar event
type CalendarEvent struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// ExchangeToken models an exchange token
type ExchangeToken struct {
	Type       string `json:"token_type"`
	ExpireTime int    `json:"expires_in"`
	Token      string `json:"access_token"`
}

// ExchangeEvent models an event returned by microsoft exchange service
type ExchangeEvent struct {
	ID                         string       `json:"id"`
	CreatedDateTime            string       `json:"createdDateTime"`
	LastModifiedDateTime       string       `json:"lastModifiedDateTime"`
	ChangeKey                  string       `json:"changeKey"`
	Categories                 string       `json:"categories"`
	OriginalStartTimeZone      string       `json:"originalStartTimeZone"`
	OriginalEndTimeZone        string       `json:"originalEndTimeZone"`
	ICalUID                    string       `json:"iCalUId"`
	ReminderMinutesBeforeStart int          `json:"reminderMinutesBeforeStart"`
	IsReminderOn               bool         `json:"isReminderOn"`
	HasAttachments             bool         `json:"hasAttachments"`
	Subject                    string       `json:"subject"`
	BodyPreview                string       `json:"bodyPreview"`
	Importance                 string       `json:"importance"`
	Sensitivity                string       `json:"sensitivity"`
	IsAllDay                   bool         `json:"isAllDay"`
	IsCancelled                bool         `json:"isCancelled"`
	IsOrganizer                bool         `json:"isOrganizer"`
	ResponseRequested          bool         `json:"responseRequested"`
	SeriesMasterID             string       `json:"seriesMasterId"`
	ShowAs                     string       `json:"showAs"`
	EventType                  string       `json:"type"`
	WebLink                    string       `json:"webLink"`
	OnlineMeetingURL           string       `json:"onlineMeetingUrl"`
	Recurrence                 string       `json:"recurrence"`
	ResponseStatus             string       `json:"responseStatus"`
	Body                       ExchangeBody `json:"body"`
	Start                      ExchangeDate `json:"start"`
	End                        ExchangeDate `json:"end"`
	Location                   string       `json:"location"`
	Locations                  string       `json:"locations"`
	Attendees                  string       `json:"attendees"`
	Organizer                  string       `json:"organizer"`
}

type ExchangeEventRequest struct {
	Subject   string             `json:"subject"`
	Body      ExchangeBody       `json:"body"`
	Start     ExchangeDate       `json:"start"`
	End       ExchangeDate       `json:"end"`
	Attendees []ExchangeAttendee `json:"attendees"`
}

type ExchangeBody struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type ExchangeDate struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

type ExchangeAttendee struct {
	EmailAddress ExchangeEmailAddress `json:"emailAddress"`
	Type         string               `json:"type"`
}

type ExchangeEmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}
