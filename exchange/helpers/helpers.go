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
func GetExchangeEvents(room string) ([]models.CalendarEvent, error) {
	//Get token
	token, err := GetToken()
	if err != nil {
		log.L.Errorf("Error getting exchange token | %v", err)
		return nil, fmt.Errorf("Error getting exchange token | %v", err)
	}
	bearerToken := "Bearer " + token

	//Todo: Identify proper calendar
	calendarID := "calendarID"

	//Get calendar events for the day
	requestURL := "https://outlook.office.com/api/v2.0/me/" + calendarID + "/calendarView"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.L.Errorf("Error creating get request to: %s | %v", requestURL, err)
		return nil, fmt.Errorf("Error creating get request to: %s | %v", requestURL, err)
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
		log.L.Errorf("Error sending http request to: %s | %v", requestURL, err)
		return nil, fmt.Errorf("Error sending http request to: %s | %v", requestURL, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Errorf("Error reading response body | %v", err)
		return nil, fmt.Errorf("Error reading response body | %v", err)
	}

	var respBody []models.ExchangeEvent
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("Error unmarshalling response body from json into exchange event object | %v", err)
		return nil, fmt.Errorf("Error unmarshalling response body from json into exchange event object | %v", err)
	}

	dateTimeLayout := "2006-01-02T15:04:05"
	var events []models.CalendarEvent
	for _, event := range respBody {
		eventStart, err := time.Parse(dateTimeLayout, event.Start.DateTime)
		if err != nil {}
		eventEnd, err := time.Parse(dateTimeLayout, event.End.DateTime)
		if err != nil {}

		location, err := time.LoadLocation(event.Start.TimeZone)
		if err != nil {}
		eventStart = eventStart.In(location)
		eventEnd = eventEnd.In(location)

		events = append(events, models.CalendarEvent{
			Title:         event.Subject,
			startTime: eventStart.Format("2006-01-02T15:04:05-07:00"),
			endTime:   eventEnd.Format("2006-01-02T15:04:05-07:00"),
		})
	}

	return events, nil
}

// SetExchangeEvent sets an event in exchange
func SetExchangeEvent(event models.CalendarEvent, room string) error {

	// Todo: Identify proper calendar
	calendarID := "calendarID"

	token, err := GetToken()
	if err != nil {
		log.L.Errorf("Error creating get request to: %s | %v", requestURL, err)
		return nil, fmt.Errorf("Error creating get request to: %s | %v", requestURL, err)
	}
	bearerToken := "Bearer " + token

	//Convert calendar event into exchange event
	eventStart := time.Date(event.StartTime)
	eventStartZone, _ := eventStart.Zone()
	eventEnd := time.Date(event.EndTime)
	eventEndZone, _ := eventEnd.Zone()

	requestBody := models.ExchangeEventRequest{
		Subject: event.Title,
		Start: models.ExchangeDate{
			DateTime: eventStart.Format("2006-01-02T15:04:05"),
			TimeZone: eventStartZone,
		},
		End: models.ExchangeDate{
			DateTime: eventEnd.Format("2006-01-02T15:04:05"),
			TimeZone: eventEndZone,
		},
	}

	//Send event request
	requestURL := "https://outlook.office.com/api/v2.0/me/calendars/" + calendarID + "/events"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.L.Errorf("Error creating post request to: %s | %v", requestURL, err)
		return nil, fmt.Errorf("Error creating post request to: %s | %v", requestURL, err)
	}

	request.Header.Add("Authorization", bearerToken)
	request.Header.Add("Content-Type", "application/json")
	Date
	resp, err := client.Do(request)
	if err != nil {
		log.L.Errorf("Error sending http request to: %s | %v", requestURL, err)calendarID
		return nil, fmt.Errorf("Error sending http request to: %s | %v", requestURL, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Errorf("Error reading response body | %v", err)
		return nil, fmt.Errorf("Error reading response body | %v", err)
	}

	var respBody string
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("Error unmarshalling response body from json into string | %v", err)
		return nil, fmt.Errorf("Error unmarshalling response body from json into string | %v", err)
	}

	return nil
}
