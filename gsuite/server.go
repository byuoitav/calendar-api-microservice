package main

import (
	"net/http"

	"github.com/byuoitav/calendar-api-microservice/gsuite/handlers"
	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
)

// func main() {
// 	helpers.AuthenticateClient()
// 	// fmt.Printf("Error: %s", err.Error())
// }

func main() {
	port := ":8034"
	router := common.NewRouter()

	// write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	// read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))
	router.GET("/events/:room", handlers.GetRoomEvents)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
