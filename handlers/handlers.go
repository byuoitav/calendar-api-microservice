package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/helpers"
	"github.com/byuoitav/calendar-api-microservice/models"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

//GetCalendarEvents handles getting the days calendar events
func GetCalendarEvents(ctx echo.Context) error {
	room := ctx.Param("room")
	log.L.Infof("Getting Calendar for room: %s", room)

	config, err := helpers.GetCouchConfig(room)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to get schedule configuration")
	}

	events, err := helpers.GetEvents(config)
	if err != nil {
		log.L.Errorf("Failed to get events for room: %s | %v", room, err)
		return ctx.JSON(http.StatusInternalServerError, "Failed to get room events")
	}

	log.L.Info("Retrieved calendar events")
	return ctx.JSON(http.StatusOK, events)
}

//SendEvent handles setting a calendar event
func SendEvent(ctx echo.Context) error {
	room := ctx.Param("room")
	log.L.Infof("Sending event to room: %s", room)

	config, err := helpers.GetCouchConfig(room)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to get schedule configuration")
	}

	var eventData models.CalendarEvent
	err = ctx.Bind(&eventData)
	if err != nil {
		log.L.Errorf("Failed to bind request body for: %s | %v", room, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to bind request body for: %s", room))
	}

	err = helpers.SetEvent(config, eventData)
	if err != nil {
		log.L.Errorf("Failed to send events to room: %s | %v", room, err)
		return ctx.JSON(http.StatusInternalServerError, "Failed to send room event")
	}

	return ctx.JSON(http.StatusOK, fmt.Sprintf("Event set successfully for room: %s", room))
}
