package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/couch"
	"github.com/byuoitav/calendar-api-microservice/helpers"
	"github.com/byuoitav/calendar-api-microservice/models"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
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

	events, err := getEvents(config)
	if err != nil {
		log.L.Errorf("Failed to get events for room: %s | %v", room, err)
		return ctx.JSON(http.StatusInternalServerError, "Failed to get room events")
	}

	log.L.Info("Retrieved calendar events")
	return ctx.JSON(http.StatusOK, events)
}

func getEvents(roomConfig structs.ScheduleConfig) ([]models.CalendarEvent, error) {
	switch roomConfig.CalendarType {
	case "Google":
		log.L.Infof("Calling G Suite microservice")
		return helpers.GetGSuite(roomConfig.Resource)
	case "Exchange":
		log.L.Infof("Calling Exchange microservice")
		return nil, fmt.Errorf("Exchange service is currently unavailable")
	default:
		return nil, fmt.Errorf("Room %s currently has no calendar type setting", roomConfig.ID)
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

	err = setEvents(config, eventData)
	if err != nil {
		log.L.Errorf("Failed to send events to room: %s | %v", room, err)
		return ctx.JSON(http.StatusInternalServerError, "Failed to send room event")
	}

	return ctx.JSON(http.StatusOK, fmt.Sprintf("Event set successfully for room: %s", room))
}

func setEvents(roomConfig structs.ScheduleConfig, event models.CalendarEvent) error {
	switch roomConfig.CalendarType {
	case "Google":
		log.L.Infof("Calling G Suite microservice")
		return helpers.SendGSuite(roomConfig.Resource, event)
	case "Exchange":
		log.L.Infof("Calling Exchange microservice")
		return fmt.Errorf("Exchange service is currently unavailable")
	default:
		return fmt.Errorf("Room %s currently has no calendar type setting", roomConfig.ID)
	}
}
