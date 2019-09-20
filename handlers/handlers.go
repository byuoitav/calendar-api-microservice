package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/couch"
	"github.com/byuoitav/calendar-api-microservice/helpers"
	"github.com/byuoitav/calendar-api-microservice/models"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

//GetCalendar ...
func GetCalendar(ctx echo.Context) error {
	room := ctx.Param("room")
	log.L.Infof("Getting Calendar for room: %s", room)

	config, err := couch.GetCouchConfig(room)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to get schedule configuration")
	}

	events, err := getEvents(config.ID, config.CalendarType)
	if err != nil {
		log.L.Errorf("Failed to get events for room: %s | %v", room, err)
		return ctx.JSON(http.StatusInternalServerError, "Failed to get room events")
	}

	log.L.Info("Retrieved calendar events")
	return ctx.JSON(http.StatusOK, events)
}

func getEvents(roomID string, service string) ([]models.CalendarEvent, error) {
	switch service {
	case "Google":
		log.L.Infof("Calling G Suite microservice")
		return helpers.GetGSuite(roomID)
	case "Exchange":
		log.L.Infof("Calling Exchange microservice")
		return nil, fmt.Errorf("Exchange service is currently unavailable")
	default:
		return nil, fmt.Errorf("Room %s currently has no calendar type setting", roomID)
	}
}

//SendEvent ...
func SendEvent(ctx echo.Context) error {
	room := ctx.Param("room")
	log.L.Infof("Sending event to room: %s", room)

	config, err := couch.GetCouchConfig(room)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to get schedule configuration")
	}

	var eventData models.CalendarEvent
	err = ctx.Bind(&eventData)
	if err != nil {
		log.L.Errorf("Failed to bind request body for: %s | %v", room, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to bind request body for: %s", room))
	}

	err = setEvents(config.ID, config.CalendarType, eventData)
	if err != nil {
		log.L.Errorf("Failed to send events to room: %s | %v", room, err)
		return ctx.JSON(http.StatusInternalServerError, "Failed to send room event")
	}

	return ctx.JSON(http.StatusOK, fmt.Sprintf("Event set successfully for room: %s", room))
}

func setEvents(roomID string, service string, event models.CalendarEvent) error {
	switch service {
	case "Google":
		log.L.Infof("Calling G Suite microservice")
		return helpers.SendGSuite(roomID, event)
	case "Exchange":
		log.L.Infof("Calling Exchange microservice")
		return fmt.Errorf("Exchange service is currently unavailable")
	default:
		return fmt.Errorf("Room %s currently has no calendar type setting", roomID)
	}
}
