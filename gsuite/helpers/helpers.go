package helpers

import (
	"fmt"

	"google.golang.org/api/calendar/v3"
)

const (
	urlPrefix = "https://www.googleapis.com/calendar/v3"
)

//CalEvent ...
type CalEvent struct {
	name      string
	startTime string
	endTime   string
}

//GetEvents ...
func GetEvents(room string, calSvc *calendar.Service) ([]CalEvent, error) {
	//find room calendar id
	calList, err := calSvc.CalendarList.List().Fields("items").Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve calendar list | %v", err)
	}

	var calID string
	for _, cal := range calList.Items {
		fmt.Printf("%v | %v\n", cal.Id, cal.Summary)
		if cal.Summary == room {
			calID = cal.Id
			break
		}
	}
	if calID == "" {
		return nil, fmt.Errorf("Room: %s does not have an assigned calendar", room)
	}
	//get days events

	eventList, err := calSvc.Events.List(calID).Fields("items(summary, start, end)").Do()
	// eventList, err := calSvc.Events.List(calID).Fields("items(summary, start, end)").TimeMin().TimeMax().Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve events | %v", err)
	}

	for _, event := range eventList.Items {
		fmt.Printf("Event %v: Start: %v, End: %v\n", event.Summary, event.Start.DateTime, event.End.DateTime)
	}

	return nil, nil
}

// listRes, err := service.CalendarList.List().Fields("items").Do()
// 	for _, v := range listRes.Items {
// 		fmt.Printf("%v | %v\n", v.Id, v.Summary)
// 	}

// 	fmt.Println("Getting event lists")

// 	for _, cal := range listRes.Items {
// 		id := cal.Id
// 		res, err := service.Events.List(id).Fields("items(updated,summary)", "summary").Do()
// 		if err != nil {
// 			fmt.Printf("Unable to retrieve calendar events list: %v", err)
// 			return err
// 		}
// 		for _, v := range res.Items {
// 			fmt.Printf("Calendar ID %q event: %v: %q\n", id, v.Updated, v.Summary)
// 		}
// 		fmt.Printf("Calendar ID %q Summary: %v\n", id, res.Summary)
// 	}

// func findRoom(room string, )
