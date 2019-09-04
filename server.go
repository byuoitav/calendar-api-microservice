package main

import (
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/handlers"
	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
)

func main() {
	port := ":8033"
	router := common.NewRouter()

	write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))

	read.GET("/calendar/:room", handlers.GetCalendar)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
