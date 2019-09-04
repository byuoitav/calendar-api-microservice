package handlers

import (
	"github.com/labstack/echo"
)

//GetCalendar ...
func GetCalendar(ctx echo.Context) error {
	room := ctx.Param("room")
	//Todo: Get calendar configuration
	//Todo: Hit appropriate calendar microservice
	//Todo: Return calendar events as json
}
