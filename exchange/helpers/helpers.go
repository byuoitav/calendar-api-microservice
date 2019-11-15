package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/byuoitav/calendar-api-microservice/exchange/models"
	"github.com/byuoitav/common/log"
)

// GetExchangeEvents gets events from exchange
func GetExchangeEvents() ([]models.CalendarEvent, error) {
	//Get token
	token, err := GetToken()
	if err != nil {
		log.L.Errorf("")
		return nil, fmt.Errorf("")
	}
	bearerToken := "Bearer " + token

	//Todo: Identify proper calendar

	//Get calendar events for the day

	requestURL := "https://outlook.office.com/api/v2.0/me/{calendarID}/calendarView"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.L.Errorf("")
		return nil, fmt.Errorf("")
	}

	request.Header.Add("Authorization", bearerToken)

	currentTime := time.Now()
	currentDayBeginning := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	currentDayEnding := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

	query := request.URL.Query()
	query.Add("startDateTime", currentDayBeginning.Format("2006-01-02T15:04:05")+".0000000")
	query.Add("endDateTime", currentDayEnding.Format("2006-01-02T15:04:05")+".0000000")
	request.URL.RawQuery = query.Encode()

	resp, err := client.Do(request)
	if err != nil {
		log.L.Errorf("")
		return nil, fmt.Errorf("")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Errorf("")
		return nil, fmt.Errorf("")
	}

	var respBody []models.ExchangeEvent
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("")
		return nil, fmt.Errorf("")
	}

	var events []models.CalendarEvent
	for _, event := range respBody {
		events = append(events, models.CalendarEvent{
			Name:          event.Subject,
			startDateTime: "",
			endDateTime:   "",
		})
	}

	return events, nil
}

// SetExchangeEvent sets an event in exchange
func SetExchangeEvent(event models.CalendarEvent) error {

	//Todo: Identify proper calendar

	token, err := GetToken()
	if err != nil {
		log.L.Errorf("")
		return fmt.Errorf("")
	}
	bearerToken := "Bearer " + token

	//Convert calendar event into exchange event
	requestBody := models.ExchangeEventRequest{
		Subject: event.Name,
		Start: models.ExchangeDate{
			DateTime: time.Date(event.StartTime).Format("2006-01-02T15:04:05"),
			TimeZone: "", // Get timezone
		},
		End: models.ExchangeDate{
			DateTime: time.Date(event.EndTime).Format("2006-01-02T15:04:05"),
			TimeZone: "", // Get timezone
		},
	}

	//Send event request

	requestURL := "https://outlook.office.com/api/v2.0/me/calendars/{calendarID}/events"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.L.Errorf("")
		return fmt.Errorf("")
	}

	request.Header.Add("Authorization", bearerToken)

	resp, err := client.Do(request)
	if err != nil {
		log.L.Errorf("")
		return fmt.Errorf("")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Errorf("")
		return fmt.Errorf("")
	}

	var respBody string
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("")
		return fmt.Errorf("")
	}

	return nil
}
