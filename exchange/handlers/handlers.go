package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/exchange/helpers"
	"github.com/byuoitav/calendar-api-microservice/exchange/models"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

// GetRoomEvents handles getting the days event for the given room
func GetRoomEvents(ctx echo.Context) error {
	roomName := ctx.Param("room")
	resource := cts.Param("resource")

	events, err := helpers.GetExchangeEvents(roomName, resource)
	if err != nil {
		log.L.Errorf("Failed to get exchange events for: %s | %v", roomName, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to get events for: %s", roomName))
	}

	log.L.Infof("Successfully retreived events for: %s", roomName)
	return ctx.JSON(http.StatusOK, events)
}

// AddRoomEvent handles adding an event to the calendar for the given room
func AddRoomEvent(ctx echo.Context) error {
	roomName := ctx.Param("room")
	resource := ctx.Param("resource")

	var eventData models.CalendarEvent
	err := ctx.Bind(&eventData)
	if err != nil {
		log.L.Errorf("Failed to bind request body for: %s | %v", roomName, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to bind request body for: %s", roomName))
	}

	err = helpers.SetExchangeEvent(eventData, roomName)
	if err != nil {
		log.L.Errorf("Failed to send exchange event | %v", err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to send exchange event for: %s", roomName))
	}

	log.L.Infof("Successfully created event for: %s", roomName)
	return ctx.JSON(http.StatusOK, fmt.Sprintf("Successfully created event for: %s", roomName))
}
