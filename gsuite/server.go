package main

import (
	"flag"
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/gsuite/handlers"
	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
)

func main() {
	port := flag.String("p", ":11001", "Port value")
	flag.Parse()
	router := common.NewRouter()

	// write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	// read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))
	router.GET("/events/:room", handlers.GetRoomEvents)
	router.PUT("/events/:room", handlers.AddRoomEvent)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           *port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
