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
	serviceURL := "http://localhost"
	switch roomConfig.CalendarType {
	case "Google":
		log.L.Info("Calling G Suite microservice")
		log.L.Infof("Requesting G Suite events for resource: %s", roomConfig.Name)
		serviceURL += ":11001/events/" + roomConfig.Resource
	case "Exchange":
		log.L.Info("Calling Exchange microservice")
		log.L.Infof("Requesting Exchange events for resource: %s", roomConfig.Name)
		serviceURL += ":11002/events/" + roomConfig.CalendarName + "/" + roomConfig.Resource
	case "TeamUp":
		log.L.Info("Calling TeamUp microservice")
		log.L.Infof("Requesting TeamUp events for resource: %s", roomConfig.Name)
		serviceURL += ":11003/events/" + roomConfig.Resource
	default:
		return nil, fmt.Errorf("Device %s currently has no calendar type setting", roomConfig.ID)
	}
	return getEventsRequest(serviceURL)
}

//SetEvent sends a calendar event to the appropriate microservice
func SetEvent(roomConfig structs.ScheduleConfig, event models.CalendarEvent) error {
	serviceURL := "http://localhost"
	switch roomConfig.CalendarType {
	case "Google":
		log.L.Info("Calling G Suite microservice")
		log.L.Infof("Sending G Suite event to calendar for resource: %s", roomConfig.Name)
		serviceURL += ":11001/events/" + roomConfig.Resource
	case "Exchange":
		log.L.Info("Calling Exchange microservice")
		serviceURL += ":11002/events/" + roomConfig.CalendarName + "/" + roomConfig.Resource
	case "TeamUp":
		log.L.Info("Calling TeamUp microservice")
		log.L.Infof("Sending TeamUp event to calendar for resource: %s", roomConfig.Name)
		serviceURL += ":11003/events/" + roomConfig.Resource
	default:
		return fmt.Errorf("Device %s currently has no calendar type setting", roomConfig.ID)
	}
	return sendEventsRequest(serviceURL, event)
}

//getEventsRequest gets all calendar events from the provided microservice
func getEventsRequest(serviceURL string) ([]models.CalendarEvent, error) {

	resp, err := http.Get(serviceURL)
	if err != nil {
		log.L.Errorf("Cannot send request to: %s | %v", serviceURL, err)
		return nil, fmt.Errorf("Cannot send request to: %s | %v", serviceURL, err)
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

	//Some services return events out of order
	events = sortEvents(events)

	//Return event array
	return events, err
}

//sendEventsRequest sends an event object to the provided microservice
func sendEventsRequest(serviceURL string, event models.CalendarEvent) error {

	requestBody, err := json.Marshal(event)
	if err != nil {
		log.L.Errorf("Cannot convert event to JSON | %v", err)
		return fmt.Errorf("Cannot convert event to JSON | %v", err)
	}

	log.L.Infof("Event object: %s", requestBody)

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, serviceURL, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json")
	if err != nil {
		log.L.Errorf("Cannot make HTTP Put request to: %s | %v", serviceURL, err)
		return fmt.Errorf("Cannot make HTTP Put request to: %s | %v", serviceURL, err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.L.Errorf("Cannot send request to: %s | %v", serviceURL, err)
		return fmt.Errorf("Cannot send request to: %s | %v", serviceURL, err)
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

func sortEvents(events []models.CalendarEvent) []models.CalendarEvent {
	quickSort(events, 0, len(events)-1)
	return events
}

func quickSort(events []models.CalendarEvent, low int, high int) {
	if low < high {
		partIndex := partition(events, low, high)
		quickSort(events, low, partIndex-1)
		quickSort(events, partIndex+1, high)
	}
}

func partition(events []models.CalendarEvent, low int, high int) int {
	pivot := events[high]
	swap1 := low
	for swap2 := low; swap2 < high; swap2++ {
		if events[swap2].StartTime < pivot.StartTime {
			swapValue(events, swap1, swap2)
			swap1++
		}
	}
	swapValue(events, swap1, high)
	return swap1
}

func swapValue(events []models.CalendarEvent, index1 int, index2 int) {
	temp := events[index1]
	events[index1] = events[index2]
	events[index2] = temp
}
