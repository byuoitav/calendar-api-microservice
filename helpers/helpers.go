package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/models"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

//GetEvents gets calendar events from the appropriate microservice
func GetEvents(roomConfig structs.ScheduleConfig) ([]models.CalendarEvent, error) {
	switch roomConfig.CalendarType {
	case "Google":
		log.L.Info("Calling G Suite microservice")
		log.L.Infof("Requesting G Suite events for resource: %s", roomConfig.Name)
		return GetGSuite(roomConfig.Resource, "http://localhost:8034/events/")
	case "Exchange":
		log.L.Info("Calling Exchange microservice")
		return nil, fmt.Errorf("Exchange service is currently unavailable")
	case "TeamUp":
		log.L.Info("Calling TeamUp microservice")
		log.L.Infof("Requesting TeamUp events for resource: %s", roomConfig.Name)
		return getEventsRequest(roomConfig.Resource, "http://localhost:8036/events/")
	default:
		return nil, fmt.Errorf("Device %s currently has no calendar type setting", roomConfig.ID)
	}
}

//SetEvent sends a calendar event to the appropriate microservice
func SetEvent(roomConfig structs.ScheduleConfig, event models.CalendarEvent) error {
	switch roomConfig.CalendarType {
	case "Google":
		log.L.Info("Calling G Suite microservice")
		log.L.Infof("Sending G Suite event to calendar for resource: %s", roomConfig.Name)
		return SendGSuite(roomConfig.Resource, event, "http://localhost:8034/events/")
	case "Exchange":
		log.L.Info("Calling Exchange microservice")
		return fmt.Errorf("Exchange service is currently unavailable")
	case "TeamUp":
		log.L.Info("Calling TeamUp microservice")
		log.L.Infof("Sending TeamUp event to calendar for resource: %s", roomConfig.Name)
		return sendEventsRequest(roomConfig.Resource, event, "http://localhost:8036/events/")
	default:
		return fmt.Errorf("Device %s currently has no calendar type setting", roomConfig.ID)
	}
}

//getEventsRequest gets all calendar events from the provided microservice
func getEventsRequest(room string, serviceURL string) ([]models.CalendarEvent, error) {
	//Call the gsuite microservice
	requestURL := serviceURL + room
	resp, err := http.Get(requestURL)
	if err != nil {
		log.L.Errorf("Cannot send request to: %s | %v", requestURL, err)
		return nil, fmt.Errorf("Cannot send request to: %s | %v", requestURL, err)
	}
	defer resp.Body.Close()

	//Translate events
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Errorf("Error resolving response body | %v", err)
		return nil, fmt.Errorf("Error resolving response body | %v", err)
	}

	var events []models.CalendarEvent
	err = json.Unmarshal([]byte(jsonData), &events)
	if err != nil {
		log.L.Errorf("Error resolving response body | %v", err)
		return nil, fmt.Errorf("Error resolving response body | %v", err)
	}

	//Return event array
	return events, err
}

//sendEventsRequest sends an event object to the provided microservice
func sendEventsRequest(room string, event models.CalendarEvent, serviceURL string) error {

	requestURL := serviceURL + room
	requestBody, err := json.Marshal(event)
	if err != nil {
		log.L.Errorf("Cannot convert event to JSON | %v", err)
		return fmt.Errorf("Cannot convert event to JSON | %v", err)
	}

	log.L.Infof("Event object: %s", requestBody)

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json")
	if err != nil {
		log.L.Errorf("Cannot make HTTP Put request to: %s | %v", requestURL, err)
		return fmt.Errorf("Cannot make HTTP Put request to: %s | %v", requestURL, err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.L.Errorf("Cannot send request to: %s | %v", requestURL, err)
		return fmt.Errorf("Cannot send request to: %s | %v", requestURL, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Errorf("Error resolving response body | %v", err)
		return fmt.Errorf("Error resolving response body | %v", err)
	}

	var respBody string
	err = json.Unmarshal([]byte(body), &respBody)
	log.L.Info(respBody)
	return nil
}
