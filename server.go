package main

import (
	"flag"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/handlers"
	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
)

func main() {
	port := flag.String("p", ":8080", "Port value")
	flag.Parse()
	router := common.NewRouter()

	// write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	// read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))

	router.GET("/calendar/:room", handlers.GetCalendarEvents)
	router.PUT("/calendar/:room", handlers.SendEvent)
	router.GET("/config", handlers.GetConfig)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	// Set up static frontend
	router.Static("/", "web-dist")

	server := http.Server{
		Addr:           *port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
