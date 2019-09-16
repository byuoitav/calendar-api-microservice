package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/models"
	"github.com/byuoitav/common/log"
)

//GetGSuite ...
func GetGSuite(room string) ([]models.GSuiteEvent, error) {
	//Todo: Figure out authentication keys, etc.
	log.L.Infof("Requesting G Suite events for room: %s", room)
	//Call the gsuite microservice
	requestURL := "http://localhost:8034/events/" + room
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

	var events []models.GSuiteEvent
	err = json.Unmarshal([]byte(jsonData), &events)
	if err != nil {
		log.L.Errorf("Error resolving response body | %v", err)
		return nil, fmt.Errorf("Error resolving response body | %v", err)
	}

	//Return event array
	return events, err
}

//SendGSuite ...
func SendGSuite(room string, event models.GSuiteEvent) error {
	log.L.Infof("Sending G Suite event to calendar for room: %s", room)

	requestURL := "http://localhost:8034/events/" + room
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
