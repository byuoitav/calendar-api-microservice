package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/byuoitav/calendar-api-microservice/teamup/models"
)

const (
	apiKey     = "TEAMUP_API_KEY"
	password   = "TEAMUP_PASSWORD"
	calendarID = "TEAMUP_CALENDAR_ID"
)

// GetTeamUpEvents sends a request to the teamup api to get all the days events for the given calendar
func GetTeamUpEvents(subCalName string) ([]models.CalendarEvent, error) {
	subCalID, err := GetSubcalendarID(subCalName)
	if err != nil {
		return nil, fmt.Errorf("Error getting subcalendar id | %s", err.Error())
	}

	client := http.Client{}
	calID := os.Getenv(calendarID)

	url := "https://api.teamup.com/" + calID + "/events"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating get request | %s", err.Error())
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Teamup-Token", os.Getenv(apiKey))
	if os.Getenv(password) != "" {
		request.Header.Set("Teamup-Password", os.Getenv(password))
	}

	query := request.URL.Query()
	query.Add("subcalendarId[]", subCalID)
	request.URL.RawQuery = query.Encode()

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error sending get request | %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body | %s", err.Error())
	}

	var eventResponse models.TeamUpEventResponse
	err = json.Unmarshal(body, &eventResponse)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling json event data | %s", err.Error())
	}

	// Convert from teamup event to calendar event
	var calendarEvents []models.CalendarEvent
	for _, event := range eventResponse.Events {
		calEvent := &models.CalendarEvent{
			Title:     event.Title,
			StartTime: event.StartDate,
			EndTime:   event.EndDate,
		}
		calendarEvents = append(calendarEvents, *calEvent)
	}

	return calendarEvents, nil
}

// SetTeamUpEvent sends a request to the teamup api to set a calendar event
func SetTeamUpEvent(subCalName string, event models.CalendarEvent) error {
	subCalIDStr, err := GetSubcalendarID(subCalName)
	if err != nil {
		return fmt.Errorf("Error getting subcalendar id | %s", err.Error())
	}
	subCalID, err := strconv.Atoi(subCalIDStr)
	if err != nil {
		return fmt.Errorf("Error converting subcalendar id to int | %s", err.Error())
	}
	// Translate event to team up event
	teamUpEvent := &models.TeamUpEventSend{
		SubCalendarID: subCalID,
		Title:         event.Title,
		StartDate:     event.StartTime,
		EndDate:       event.EndTime,
		AllDay:        false,
	}

	client := http.Client{}
	calID := os.Getenv(calendarID)

	requestBody, err := json.Marshal(teamUpEvent)
	if err != nil {
		return fmt.Errorf("Error marshalling event data into json | %s", err.Error())
	}

	url := "https://api.teamup.com/" + calID + "/events"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("Error creating post request | %s", err.Error())
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Teamup-Token", os.Getenv(apiKey))
	if os.Getenv(password) != "" {
		request.Header.Set("Teamup-Password", os.Getenv(password))
	}

	// Send request
	resp, err := client.Do(request)
	if err != nil {
		// Error
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Error
		return err
	}

	// Get 201 response
	if resp.StatusCode != 201 {
		return fmt.Errorf("Incorrect response code | %s", body)
	}
	return nil
}

// GetSubcalendarID sends a request to the teamup api to get the appropriate subcalendar id
func GetSubcalendarID(subcalendarName string) (string, error) {
	calID := os.Getenv(calendarID)
	client := http.Client{}

	url := "https://api.teamup.com/" + calID + "/subcalendars"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// Error
		return "", err
	}

	apikey := os.Getenv(apiKey)
	request.Header.Set("Teamup-Token", apikey)
	if os.Getenv(password) != "" {
		request.Header.Set("Teamup-Password", os.Getenv(password))
	}

	resp, err := client.Do(request)
	if err != nil {
		// Error
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Error
		return "", err
	}

	var subcalResponse models.TeamUpSubcalendarList
	err = json.Unmarshal(body, &subcalResponse)
	if err != nil {
		// Error
		return "", fmt.Errorf("Error unmarshalling json | %v", err)
	}

	for _, subCal := range subcalResponse.Subcalendars {
		if subCal.Name == subcalendarName {
			return strconv.Itoa(subCal.ID), nil
		}
	}
	return "", nil
}
