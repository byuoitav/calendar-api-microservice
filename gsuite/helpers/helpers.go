package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/byuoitav/calendar-api-microservice/gsuite/models"
	"github.com/byuoitav/common/log"
	"google.golang.org/api/calendar/v3"
)

const (
	urlPrefix   = "https://www.googleapis.com/calendar/v3"
	userEmail   = "G_SUITE_EMAIL"
	credentials = "G_SUITE_CREDENTIALS"
)

//GetEvents finds the appropriate calendar and returns the days events
func GetEvents(room string) ([]models.CalendarEvent, error) {
	log.L.Infof("Getting events for resource: %s", room)

	calSvc, err := AuthenticateClient(os.Getenv(credentials), os.Getenv(userEmail))
	if err != nil {
		log.L.Error("Cannot authenticate client")
		return nil, fmt.Errorf("Cannot authenticate client | %v", err)
	}

	calID, err := findCalendarID(room, calSvc)
	if err != nil {
		log.L.Error("Cannot find calendar ID")
		return nil, err
	}

	currentTime := time.Now()
	currentDayBeginning := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	currentDayEnding := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

	eventList, err := calSvc.Events.List(calID).Fields("items(summary, start, end)").TimeMin(currentDayBeginning.Format("2006-01-02T15:04:05-07:00")).TimeMax(currentDayEnding.Format("2006-01-02T15:04:05-07:00")).Do()
	if err != nil {
		log.L.Errorf("Unable to retrieve events")
		return nil, fmt.Errorf("Unable to retrieve events | %v", err)
	}

	var events []models.CalendarEvent
	for _, event := range eventList.Items {
		events = append(events, models.CalendarEvent{
			Name:      event.Summary,
			StartTime: event.Start.DateTime,
			EndTime:   event.End.DateTime})
	}

	return events, err
}

//SetEvent finds the appropriate calendar and adds the given event
func SetEvent(room string, event models.CalendarEvent) error {
	calSvc, err := AuthenticateClient(os.Getenv(credentials), os.Getenv(userEmail))
	if err != nil {
		log.L.Error("Cannot authenticate client")
		return fmt.Errorf("Cannot authenticate client | %v", err)
	}

	calID, err := findCalendarID(room, calSvc)
	if err != nil {
		log.L.Error("Cannot find calendar ID")
		return err
	}

	//Translate event into g suite calendar event object
	newEvent := &calendar.Event{
		Summary: event.Name,
		Start: &calendar.EventDateTime{
			DateTime: event.StartTime,
		},
		End: &calendar.EventDateTime{
			DateTime: event.EndTime,
		},
	}

	newEvent, err = calSvc.Events.Insert(calID, newEvent).Do()
	if err != nil {
		log.L.Errorf("Unable to create event")
		return fmt.Errorf("Unable to create event | %v", err)
	}

	return nil
}

//Finds and returns calendar id based on calendar/room name
func findCalendarID(room string, calSvc *calendar.Service) (string, error) {
	calList, err := calSvc.CalendarList.List().Fields("items").Do()
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve calendar list | %v", err)
	}

	log.L.Debug("Calendar Names:")
	for _, cal := range calList.Items {
		log.L.Debugf("%s", cal.Summary)
		if cal.Summary == room {
			return cal.Id, nil
		}
	}
	return "", fmt.Errorf("Room: %s does not have an assigned calendar", room)
}
