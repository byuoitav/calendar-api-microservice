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
