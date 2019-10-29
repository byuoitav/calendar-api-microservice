package helpers

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/byuoitav/calendar-api-microservice/teamup/models"
	"google.golang.org/api/transport/http"
)

const API_KEY = "TEAM_UP_API_KEY"

// GetTeamUpEvents ...
func GetTeamUpEvents(calId string, subCalId string) (models.TeamUpEventResponse, error) {

	client := http.NewClient

	url := "https://api.teamup.com/" + calId + "/events"
	request, err := http.NewRequest("GET", url)
	if err != nil {
		// Error
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Teamup-Token", os.Getenv(API_KEY))
	query := request.URL.Query()
	query.Add("subcalendarId[]", subCalId)
	request.URL.RawQuery = query.Encode()

	resp, err := client.Do(request)
	if err != nil {
		// Error
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Error
		return nil, err
	}

	var eventResponse models.TeamUpEventResponse
	err := json.Unmarshal(body, &eventResponse)
	if err != nil {
		// Error
		return nil, err
	}

	return eventResponse, nil
}

func SetTeamUpEvent(calId string, subCalId string, event models.CalendarEvent) error {

	client := http.NewClient

	url := "https://api.teamup.com/" + calId + "/events"
	request, err := http.NewRequest("POST", url)
	if err != nil {
		return err
	}
	
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Teamup-Token", os.Getenv(API_KEY))

	// Translate event to team up event
	teamUpEvent := models.TeamUpEventSend{
		SubCalendarID = subCalId,
		Title = event.Name,
		StartDate = event.StartTime,
		EndDate = event.EndTime,
		AllDay = false
	}

	// Send request
	resp, err := client.Do(request)
	if err != nil {
		// Error
		return err
	}
	defer resp.Body.Close()

	// Get 201 response
	if resp.StatusCode != 201 {
		// Error - Provide body?
		return fmt.Errorf("Incorrect response code | ")
	}
	return nil
}
