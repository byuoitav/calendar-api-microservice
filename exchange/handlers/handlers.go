package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/exchange/models"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

// GetRoomEvents handles getting the days event for the given room
func GetRoomEvents(ctx *echo.Context) error {
	roomName := ctx.param("room")

	return nil
}

// AddRoomEvent handles adding an event to the calendar for the given room
func AddRoomEvent(ctx *echo.Context) error {
	roomName := ctx.param("room")
	var eventData models.CalendarEvent
	err := ctx.Bind(&eventData)
	if err != nil {
		log.L.Errorf("Failed to bind request body for: %s | %v", roomName, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to bind request body for: %s", roomName))
	}

	return nil
}
