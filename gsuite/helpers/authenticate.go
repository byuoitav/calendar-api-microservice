package helpers

import (
	"io/ioutil"
	"os"

	"github.com/byuoitav/common/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

//AuthenticateClient ...
func AuthenticateClient(credentials string) (*calendar.Service, error) {
	pwd, _ := os.Getwd()
	data, err := ioutil.ReadFile(pwd + credentials)
	if err != nil {
		log.L.Errorf("Can't read project key file | %s", err.Error())
		return nil, err
	}

	log.L.Infof("Signing JWT")
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.L.Errorf("Can't sign JWT | %s", err.Error())
		return nil, err
	}

	log.L.Infof("Getting authorization")
	client := conf.Client(oauth2.NoContext)

	return calendar.New(client)
}

//Scope: https://www.googleapis.com/auth/calendar
