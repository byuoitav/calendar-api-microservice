package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/gsuite/models"

	"github.com/byuoitav/calendar-api-microservice/gsuite/helpers"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

//GetRoomEvents ...
func GetRoomEvents(ctx echo.Context) error {
	roomName := ctx.Param("room")
	calendarService, err := helpers.AuthenticateClient("/helpers/go-calendar.json")
	if err != nil {
		log.L.Errorf("Cannot authenticate client: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, "Cannot authenticate client")
	}

	roomEvents, err := helpers.GetEvents(roomName, calendarService)
	if err != nil {
		log.L.Errorf("Error getting events: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, "Cannot get events")
	}

	return ctx.JSON(http.StatusOK, roomEvents)
}

//AddRoomEvent ...
func AddRoomEvent(ctx echo.Context) error {
	roomName := ctx.Param("room")
	var eventData models.CalendarEvent
	err := ctx.Bind(&eventData)
	if err != nil {
		log.L.Errorf("Failed to bind request body for: %s | %v", roomName, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to bind request body for: %s", roomName))
	}

	calendarService, err := helpers.AuthenticateClient("/helpers/go-calendar.json")
	if err != nil {
		log.L.Errorf("Cannot authenticate client: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, "Cannot authenticate client")
	}

	err = helpers.SetEvent(roomName, eventData, calendarService)
	if err != nil {
		log.L.Errorf("Failed to set event for room: %s | %v", roomName, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to set event for: %s", roomName))
	}

	return ctx.JSON(http.StatusOK, fmt.Sprintf("Event set for room: %s", roomName))
}
