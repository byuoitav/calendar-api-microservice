package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/gsuite/models"

	"github.com/byuoitav/calendar-api-microservice/gsuite/helpers"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

//GetRoomEvents handles getting the events for a given room
func GetRoomEvents(ctx echo.Context) error {
	roomName := ctx.Param("room")

	roomEvents, err := helpers.GetEvents(roomName)
	if err != nil {
		log.L.Errorf("Error getting events | %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get events for room: %s", roomName))
	}

	return ctx.JSON(http.StatusOK, roomEvents)
}

//AddRoomEvent handles adding an event for a given room
func AddRoomEvent(ctx echo.Context) error {
	roomName := ctx.Param("room")
	var eventData models.CalendarEvent
	err := ctx.Bind(&eventData)
	if err != nil {
		log.L.Errorf("Failed to bind request body for: %s | %v", roomName, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to bind request body for: %s", roomName))
	}

	err = helpers.SetEvent(roomName, eventData)
	if err != nil {
		log.L.Errorf("Failed to set event for room: %s | %v", roomName, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to set event for: %s", roomName))
	}

	return ctx.JSON(http.StatusOK, fmt.Sprintf("Event set for room: %s", roomName))
}
