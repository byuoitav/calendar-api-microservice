package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/helpers"
	"github.com/byuoitav/calendar-api-microservice/models"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

//GetCalendar ...
func GetCalendar(ctx echo.Context) error {
	room := ctx.Param("room")
	log.L.Infof("Getting Calendar for room: %s", room)

	//Todo: Get calendar configuration
	//Todo: Hit appropriate calendar microservice

	events, err := helpers.GetGSuite(room)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to get G Suite events")
	}

	log.L.Info("Retrieved calendar events")
	return ctx.JSON(http.StatusOK, events)
}

//SendEvent ...
func SendEvent(ctx echo.Context) error {
	room := ctx.Param("room")
	log.L.Infof("Sending Event to calendar for room: %s", room)

	//Todo: Get calendar configuration
	//Todo: Hit appropriate calendar microservice

	var eventData models.GSuiteEvent
	err := ctx.Bind(&eventData)
	if err != nil {
		log.L.Errorf("Failed to bind request body for: %s | %v", room, err)
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to bind request body for: %s", room))
	}

	err = helpers.SendGSuite(room, eventData)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to send G Suite event")
	}

	return ctx.JSON(http.StatusOK, fmt.Sprintf("Event set successfully for room: %s", room))
}
