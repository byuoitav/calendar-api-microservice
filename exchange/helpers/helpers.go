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
func GetExchangeEvents(room string, resource string) ([]models.CalendarEvent, error) {
	token, err := GetToken()
	if err != nil {
		log.L.Errorf("Error getting exchange token | %v", err)
		return nil, fmt.Errorf("Error getting exchange token | %v", err)
	}
	bearerToken := "Bearer " + token

	calendarID, err := getCalendarID(room, resource, token)
	if err != nil {
		log.L.Errorf("Error getting calendar ID for room: %s | %v", room, err)
		return nil, fmt.Errorf("Error getting calendar ID for room: %s | %v", room, err)
	}

	//Get calendar events for the day
	requestURL := "https://outlook.office.com/api/v2.0/users/" + resource + "/calendars/" + calendarID + "/calendarView"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.L.Errorf("Error creating get request to: %s | %v", requestURL, err)
		return nil, fmt.Errorf("Error creating get request to: %s | %v", requestURL, err)
	}

	request.Header.Add("Authorization", bearerToken)

	loc, _ := time.LoadLocation("UTC")
	currentTime := time.Now()
	currentDayBeginning := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	currentDayEnding := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())
	currentDayBeginning = currentDayBeginning.In(loc)
	currentDayEnding = currentDayEnding.In(loc)

	query := request.URL.Query()
	query.Add("startDateTime", currentDayBeginning.Format("2006-01-02T15:04:05"))
	query.Add("endDateTime", currentDayEnding.Format("2006-01-02T15:04:05"))
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

	var respBody models.ExchangeEventResponse
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("Error unmarshalling response body from json into exchange event object | %v", err)
		return nil, fmt.Errorf("Error unmarshalling response body from json into exchange event object | %v", err)
	}

	dateTimeLayout := "2006-01-02T15:04:05"
	var events []models.CalendarEvent
	for _, event := range respBody.Events {
		eventStart, err := time.Parse(dateTimeLayout, event.Start.DateTime)
		if err != nil {
		}
		eventEnd, err := time.Parse(dateTimeLayout, event.End.DateTime)
		if err != nil {
		}

		timeZone, _ := time.Now().Zone()
		location, _ := time.LoadLocation(timeZone)
		eventStart = eventStart.In(location)
		eventEnd = eventEnd.In(location)

		events = append(events, models.CalendarEvent{
			Title:     event.Subject,
			StartTime: eventStart.Format("2006-01-02T15:04:05-07:00"),
			EndTime:   eventEnd.Format("2006-01-02T15:04:05-07:00"),
		})
	}

	return events, nil
}

// SetExchangeEvent sets an event in exchange
func SetExchangeEvent(event models.CalendarEvent, room string, resource string) error {

	token, err := GetToken()
	if err != nil {
		log.L.Errorf("Error getting auth token | %v", err)
		return fmt.Errorf("Error getting auth token | %v", err)
	}
	bearerToken := "Bearer " + token

	calendarID, err := getCalendarID(room, resource, token)
	if err != nil {
		log.L.Errorf("Error getting calendar ID for room: %s | %v", room, err)
		return fmt.Errorf("Error getting calendar ID for room: %s | %v", room, err)
	}

	//Convert calendar event into exchange event
	loc, _ := time.LoadLocation("UTC")
	dateTimeLayout := "2006-01-02T15:04:05-07:00"

	eventStart, err := time.Parse(dateTimeLayout, event.StartTime)
	if err != nil {
	}
	eventStart = eventStart.In(loc)
	// eventStartZone, _ := eventStart.Zone()

	eventEnd, err := time.Parse(dateTimeLayout, event.EndTime)
	if err != nil {
	}
	eventEnd = eventEnd.In(loc)
	// eventEndZone, _ := eventEnd.Zone()

	requestBodyStruct := models.ExchangeEventRequest{ // Will probably need to add attendees and body
		Subject: event.Title,
		Start: models.ExchangeDate{
			DateTime: eventStart.Format("2006-01-02T15:04:05"),
			TimeZone: "Etc/GMT", // Hard code UTC
		},
		End: models.ExchangeDate{
			DateTime: eventEnd.Format("2006-01-02T15:04:05"),
			TimeZone: "Etc/GMT", // Hard code UTC
		},
	}
	requestBodyString, err := json.Marshal(requestBodyStruct)
	if err != nil {
		log.L.Errorf("Error converting request body to json string | %v", err)
		return fmt.Errorf("Error converting request body to json string | %v", err)
	}

	requestURL := "https://outlook.office.com/api/v2.0/users/" + resource + "/calendars/" + calendarID + "/events"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestBodyString))
	if err != nil {
		log.L.Errorf("Error creating post request to: %s | %v", requestURL, err)
		return fmt.Errorf("Error creating post request to: %s | %v", requestURL, err)
	}

	request.Header.Add("Authorization", bearerToken)
	request.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		log.L.Errorf("Error sending http request to: %s | %v", requestURL, err)
		return fmt.Errorf("Error sending http request to: %s | %v", requestURL, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var respBody string
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("Error unmarshalling response body from json into string | %v", err)
		return fmt.Errorf("Error unmarshalling response body from json into string | %v", err)
	}

	return nil
}

func getCalendarID(calendarName string, resource string, token string) (string, error) {

	requestURL := "https://outlook.office.com/api/v2.0/users/" + resource + "/calendars"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.L.Errorf("Error creating get request to: %s | %v", requestURL, err)
		return "", fmt.Errorf("Error creating get request to: %s | %v", requestURL, err)
	}

	request.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		log.L.Errorf("Error sending http request to: %s | %v", requestURL, err)
		return "", fmt.Errorf("Error sending http request to: %s | %v", requestURL, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Errorf("Error reading response body | %v", err)
		return "", fmt.Errorf("Error reading response body | %v", err)
	}

	var respBody models.ExchangeCalenderResponse
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("Error unmarshalling response body from json into exchange calendar object | %v", err)
		return "", fmt.Errorf("Error unmarshalling response body from json into exchange calendar object | %v", err)
	}

	// Locate the proper calendar
	if len(respBody.Calendars) > 1 {
		for _, calendar := range respBody.Calendars {
			if calendarName == calendar.Name {
				return calendar.ID, nil
			}
		}
	}
	return respBody.Calendars[0].ID, nil
}
