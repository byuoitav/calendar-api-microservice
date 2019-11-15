package helpers

import (
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/exchange/models"
)

// GetExchangeEvents gets events from exchange
func GetExchangeEvents() ([]models.CalendarEvent, error) {
	//Identify calendar
	//Get calendar events for the day

	requestURL := "https://graph.microsoft.com/me/calendarView"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)

	query := request.URL.Query()
	query.Add("startDateTime", "2019-11-15T00:00:00.0000000")
	query.Add("endDateTime", "2019-11-15T23:59:59.0000000")
	request.URL.RawQuery = query.Encode()

	//Unmarshal into a list of events
	return nil, nil
}

// SetExchangeEvent sets an event in exchange
func SetExchangeEvent(event models.CalendarEvent) error {
	//Identify calendar
	//Convert calendar event into exchange event
	//Send event request
	return nil
}
