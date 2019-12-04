package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/byuoitav/calendar-api-microservice/exchange/models"
	"github.com/byuoitav/common/log"
)

const (
	clientID     = "AZURE_AD_CLIENT_ID"
	clientSecret = "AZURE_AD_CLIENT_SECRET"
	tennantID    = "AZURE_AD_TENNANT_ID"
)

// GetToken sends a request to microsoft to get a bearer token and returns the result
func GetToken() (string, error) {

	bodyParams := url.Values{}
	bodyParams.Set("client_id", os.Getenv(clientID))
	bodyParams.Set("scope", "https://outlook.office.com/.default")
	bodyParams.Set("client_secret", os.Getenv(clientSecret))
	bodyParams.Set("grant_type", "client_credentials")

	requestURL := "https://login.microsoftonline.com/" + os.Getenv(tennantID) + "/oauth2/v2.0/token"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(bodyParams.Encode()))
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")
	if err != nil {
		log.L.Errorf("Cannot make HTTP Post request to: %s | %v", requestURL, err)
		return "", fmt.Errorf("Cannot make HTTP Post request to: %s | %v", requestURL, err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.L.Errorf("Cannot send request to: %s | %v", requestURL, err)
		return "", fmt.Errorf("Cannot send request to: %s | %v", requestURL, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Errorf("Error resolving response body | %v", err)
		return "", fmt.Errorf("Error resolving response body | %v", err)
	}

	var respBody models.ExchangeToken
	err = json.Unmarshal([]byte(body), &respBody)
	if err != nil {
		log.L.Errorf("Error unmarshalling json body | %v", err)
		return "", fmt.Errorf("Error unmarshalling json body | %v", err)
	}

	return respBody.Token, nil
}
