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
	log.L.Info("Getting exchange events...")
	log.L.Info("Getting auth token...")
	token, err := GetToken()
	if err != nil {
		log.L.Errorf("Error getting exchange token | %v", err)
		return nil, fmt.Errorf("Error getting exchange token | %v", err)
	}
	bearerToken := "Bearer " + token

	log.L.Info("Getting calendarID...")
	calendarID, err := getCalendarID(room, resource, token)
	if err != nil {
		log.L.Errorf("Error getting calendar ID for room: %s | %v", room, err)
		return nil, fmt.Errorf("Error getting calendar ID for room: %s | %v", room, err)
	}

	log.L.Info("Preparing http request for exchange events...")
	requestURL := "https://outlook.office.com/api/v2.0/users/" + resource + "/calendars/" + calendarID + "/calendarView"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.L.Errorf("Error creating get request to: %s | %v", requestURL, err)
		return nil, fmt.Errorf("Error creating get request to: %s | %v", requestURL, err)
	}

	request.Header.Add("Authorization", bearerToken)

	log.L.Info("Prepping http request query parameters...")
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

	log.L.Info("Sending http request for exchange events...")
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

	log.L.Info("Checking for proper response...")
	var respBody models.ExchangeEventResponse
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("Error unmarshalling response body from json into exchange event object | %v", err)
		return nil, fmt.Errorf("Error unmarshalling response body from json into exchange event object | %v", err)
	}

	log.L.Info("Converting exchange events into usable form...")
	dateTimeLayout := "2006-01-02T15:04:05"
	var events []models.CalendarEvent
	for _, event := range respBody.Events {
		eventStart, err := time.Parse(dateTimeLayout, event.Start.DateTime)
		if err != nil {
			log.L.Errorf("Error parsing exchange event start time into go time struct | %v", err)
			return fmt.Errorf("Error parsing exchange event start time into go time struct | %v", err)
		}
		eventEnd, err := time.Parse(dateTimeLayout, event.End.DateTime)
		if err != nil {
			log.L.Errorf("Error parsing exchange event end time into go time struct | %v", err)
			return fmt.Errorf("Error parsing exchange event end time into go time struct | %v", err)
		}

		timeZone, _ := time.Now().Zone()
		location, _ := time.LoadLocation(timeZone)log.L.Info()
		events = append(events, models.CalendarEvent{
			Title:     event.Subject,
			StartTime: eventStart.Format("2006-01-02T15:04:05-07:00"),
			EndTime:   eventEnd.Format("2006-01-02T15:04:05-07:00"),
		})
	}

	log.L.Info("Successfully retrieved exchange events...")
	return events, nil
}

// SetExchangeEvent sets an event in exchange
func SetExchangeEvent(event models.CalendarEvent, room string, resource string) error {
	log.L.Info("Sending exchange event...")
	log.L.Info("Getting auth token...")
	token, err := GetToken()
	if err != nil {
		log.L.Errorf("Error getting auth token | %v", err)
		return fmt.Errorf("Error getting auth token | %v", err)
	}
	bearerToken := "Bearer " + token

	log.L.Info("Getting calendarID...")
	calendarID, err := getCalendarID(room, resource, token)
	if err != nil {
		log.L.Errorf("Error getting calendar ID for room: %s | %v", room, err)
		return fmt.Errorf("Error getting calendar ID for room: %s | %v", room, err)
	}

	//Convert calendar event into exchange event
	log.L.Info("Converting calendar event into exchange format...")
	loc, _ := time.LoadLocation("UTC")
	dateTimeLayout := "2006-01-02T15:04:05-07:00"

	eventStart, err := time.Parse(dateTimeLayout, event.StartTime)
	if err != nil {
		log.L.Errorf("Error parsing event start time into go time struct | %v", err)
		return fmt.Errorf("Error parsing event start time into go time struct | %v", err)
	}
	eventStart = eventStart.In(loc)

	eventEnd, err := time.Parse(dateTimeLayout, event.EndTime)
	if err != nil {
		log.L.Errorf("Error parsing event start time into go time struct | %v", err)
		return fmt.Errorf("Error parsing event start time into go time struct | %v", err)
	}
	eventEnd = eventEnd.In(loc)

	log.L.Info("Creating exchange event http post request body...")
	requestBodyStruct := models.ExchangeEventRequest{ // Will probably need to add attendees and body
		Subject: event.Title,
		Body: models.ExchangeBody{
			ContentType: "HTML",
			Content:     "",
		},
		Start: models.ExchangeDate{
			DateTime: eventStart.Format("2006-01-02T15:04:05"),
			TimeZone: "Etc/GMT", // Hard code UTC
		},
		End: models.ExchangeDate{
			DateTime: eventEnd.Format("2006-01-02T15:04:05"),
			TimeZone: "Etc/GMT", // Hard code UTC
		},
		Attendees: make([]models.ExchangeAttendee, 0),
	}
	requestBodyString, err := json.Marshal(requestBodyStruct)
	if err != nil {
		log.L.Errorf("Error converting request body to json string | %v", err)
		return fmt.Errorf("Error converting request body to json string | %v", err)
	}

	log.L.Info("Prepping exchange event http post request...")
	requestURL := "https://outlook.office.com/api/v2.0/users/" + resource + "/calendars/" + calendarID + "/events"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestBodyString))
	if err != nil {
		log.L.Errorf("Error creating post request to: %s | %v", requestURL, err)
		return fmt.Errorf("Error creating post request to: %s | %v", requestURL, err)
	}

	request.Header.Add("Authorization", bearerToken)
	request.Header.Add("Content-Type", "application/json")

	log.L.Info("Sending exchange event http post request...")
	resp, err := client.Do(request)
	if err != nil {
		log.L.Errorf("Error sending http request to: %s | %v", requestURL, err)
		return fmt.Errorf("Error sending http request to: %s | %v", requestURL, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// Report that the event was sent successfully?

	log.L.Info("Validating http response...")
	var respBody models.ExchangeEvent
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("Error sending exchange event - event not registered properly: %s | %v", string(body), err)
		return fmt.Errorf("Error sending exchange event - event not registered properly | %v", err)
	}

	log.L.Info("Event sent and registered correctly...")
	return nil
}

func getCalendarID(calendarName string, resource string, token string) (string, error) {
	log.L.Info("Prepping calenderID request...")
	requestURL := "https://outlook.office.com/api/v2.0/users/" + resource + "/calendars"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.L.Errorf("Error creating get request to: %s | %v", requestURL, err)
		return "", fmt.Errorf("Error creating get request to: %s | %v", requestURL, err)
	}

	request.Header.Add("Authorization", "Bearer "+token)

	log.L.Info("Sending calendar http get request...")
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

	log.L.Info("Validating http response...")
	var respBody models.ExchangeCalenderResponse
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.L.Errorf("Error unmarshalling response body from json into exchange calendar object | %v", err)
		return "", fmt.Errorf("Error unmarshalling response body from json into exchange calendar object | %v", err)
	}

	// Locate the proper calendar
	if len(respBody.Calendars) > 1 {
		log.L.Info("Several calendars returned, finding proper calendar...")
		for _, calendar := range respBody.Calendars {
			if calendarName == calendar.Name {
				return calendar.ID, nil
			}
		}
		log.L.Info("Given calendar name did not match. Returning first calendarID by default...")
	}
	log.L.Info("Returning calendarID...")
	return respBody.Calendars[0].ID, nil
}
