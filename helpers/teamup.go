package helpers

import "github.com/byuoitav/calendar-api-microservice/models"

// GetTeamUp retrieves events from a calendar through the TeamUp microservice
func GetTeamUp(room string) ([]models.CalendarEvent, error) {
	return nil, nil
}

// SendTeamUp sends events to a calendar through the TeamUp microservice
func SendTeamUp(room string, event models.CalendarEvent) error {
	return nil
}
