package helpers

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

//AuthenticateClient ...
func AuthenticateClient() (*Service, error) {
	pwd, _ := os.Getwd()
	data, err := ioutil.ReadFile(pwd + "/helpers/go-calendar.json")
	if err != nil {
		fmt.Printf("Can't read project key file | %s", err.Error())
		return nil, err
	}

	fmt.Println("Signing JWT")
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		fmt.Printf("Can't sign JWT | %s", err.Error())
		return nil, err
	}

	fmt.Println("Getting authorization")
	client := conf.Client(oauth2.NoContext)

	service, err := calendar.New(client)
	if err != nil {
		return nil, err
	}
	return service, err
}

//Scope: https://www.googleapis.com/auth/calendar
