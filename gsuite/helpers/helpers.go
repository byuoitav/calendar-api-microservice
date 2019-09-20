package helpers

import (
	"fmt"
	"time"

	"github.com/byuoitav/calendar-api-microservice/gsuite/models"
	"github.com/byuoitav/common/log"
	"google.golang.org/api/calendar/v3"
)

const (
	urlPrefix = "https://www.googleapis.com/calendar/v3"
)

//GetEvents ...
func GetEvents(room string, calSvc *calendar.Service) ([]models.CalendarEvent, error) {
	log.L.Infof("Getting events for room: %s", room)
	calID, err := findCalendarID(room, calSvc)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	currentDayBeginning := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	currentDayEnding := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

	eventList, err := calSvc.Events.List(calID).Fields("items(summary, start, end)").TimeMin(currentDayBeginning.Format("2006-01-02T15:04:05-07:00")).TimeMax(currentDayEnding.Format("2006-01-02T15:04:05-07:00")).Do()
	if err != nil {
		log.L.Errorf("Unable to retrieve events | %v", err)
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

//SetEvent ...
func SetEvent(room string, event models.CalendarEvent, calSvc *calendar.Service) error {
	calID, err := findCalendarID(room, calSvc)
	if err != nil {
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
		log.L.Errorf("Unable to create event | %v", err)
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

	for _, cal := range calList.Items {
		if cal.Summary == room {
			return cal.Id, nil
		}
	}
	return "", fmt.Errorf("Room: %s does not have an assigned calendar", room)
}
