package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/gsuite/helpers"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

//GetRoomEvents ...
func GetRoomEvents(ctx echo.Context) error {
	roomName := ctx.Param("room")
	calendarService, err := helpers.AuthenticateClient()
	if err != nil {
		log.L.Errorf("Cannot authenticate client: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot authenticate client: %s", err.Error()))
	}

	return ctx.JSON(http.StatusOK, "Room Events")
}
