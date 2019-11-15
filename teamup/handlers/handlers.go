package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/calendar-api-microservice/teamup/helpers"
	"github.com/byuoitav/calendar-api-microservice/teamup/models"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

const calID = "TEAMUP_CALENDAR_ID"

// GetRoomEvents handles getting events from the teamup calendar
func GetRoomEvents(ctx echo.Context) error {
	//get room
	roomName := ctx.Param("room")
	//get calendar id
	calendarID := os.Getenv(calID)

	events, err := helpers.GetTeamUpEvents(calendarID, roomName)
	if err != nil {
		log.L.Errorf("Error getting events | %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, "Cannot get events")
	}

	return ctx.JSON(http.StatusOK, events)
}

// AddRoomEvent handles adding an event to the teamup calendar
func AddRoomEvent(ctx echo.Context) error {
	//get room
	roomName := ctx.Param("room")
	//get calendar id
	calendarID := os.Getenv(calID)

	//Read body into calendar event object
	var eventData models.CalendarEvent
	err := ctx.Bind(&eventData)
	if err != nil {
		log.L.Errorf("Failed to bind request body for: %s | %v", roomName, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to bind request body for: %s", roomName))
	}

	err = helpers.SetTeamUpEvent(calendarID, roomName, eventData)
	if err != nil {
		log.L.Errorf("Error setting an event | %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, "Cannot set event")
	}

	return ctx.JSON(http.StatusOK, "Event set")
}
