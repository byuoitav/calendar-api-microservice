package helpers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/byuoitav/common/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

//AuthenticateClient ...
func AuthenticateClient(credentials string, userEmail string) (*calendar.Service, error) {
	// func AuthenticateClient(credentials string) (*calendar.Service, error) {
	ctx := context.Background()

	pwd, _ := os.Getwd()
	data, err := ioutil.ReadFile(pwd + credentials)
	if err != nil {
		log.L.Errorf("Can't read project key file | %s", err.Error())
		return nil, err
	}

	log.L.Info("Signing JWT")
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.L.Errorf("Can't sign JWT | %s", err.Error())
		return nil, err
	}
	conf.Subject = userEmail

	ts := conf.TokenSource(ctx)

	log.L.Info("Getting authorization")
	service, err := calendar.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("Can't make calendar service | %s", err.Error())
	}

	return service, nil

	// client := conf.Client(oauth2.NoContext)

	// return calendar.New(client)
}
