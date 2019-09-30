package helpers

import (
	"fmt"

	"github.com/byuoitav/calendar-api-microservice/models"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

//GetEvents gets calendar events from the appropriate microservice
func GetEvents(roomConfig structs.ScheduleConfig) ([]models.CalendarEvent, error) {
	switch roomConfig.CalendarType {
	case "Google":
		log.L.Infof("Calling G Suite microservice")
		return GetGSuite(roomConfig.Resource)
	case "Exchange":
		log.L.Infof("Calling Exchange microservice")
		return nil, fmt.Errorf("Exchange service is currently unavailable")
	default:
		return nil, fmt.Errorf("Room %s currently has no calendar type setting", roomConfig.ID)
	}
}

//SetEvent sends a calendar event to the appropriate microservice
func SetEvent(roomConfig structs.ScheduleConfig, event models.CalendarEvent) error {
	switch roomConfig.CalendarType {
	case "Google":
		log.L.Infof("Calling G Suite microservice")
		return SendGSuite(roomConfig.Resource, event)
	case "Exchange":
		log.L.Infof("Calling Exchange microservice")
		return fmt.Errorf("Exchange service is currently unavailable")
	default:
		return fmt.Errorf("Room %s currently has no calendar type setting", roomConfig.ID)
	}
}
