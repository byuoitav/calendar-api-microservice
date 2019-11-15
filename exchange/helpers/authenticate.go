package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/byuoitav/calendar-api-microservice/exchange/models"
)

const (
	// PROXY_USERNAME = "EXCHANGE_PROXY_USERNAME"
	// PROXY_PASSWORD = "EXCHANGE_PROXY_PASSWORD"
	clientID     = "AZURE_AD_CLIENT_ID"
	clientSecret = "AZURE_AD_CLIENT_SECRET"
	tennantID    = "AZURE_AD_TENNANT_ID"
)

// GetToken sends a request to microsoft to get a bearer token and returns the result
func GetToken() (string, error) {

	bodyParams := url.Values{}
	bodyParams.Set("client_id", os.Getenv(clientID))
	bodyParams.Set("scope", "https://graph.microsoft.com/.default")
	bodyParams.Set("client_secret", os.Getenv(clientSecret))
	bodyParams.Set("grant_type", "client_credentials")

	requestURL := "https://login.microsoftonline.com/" + os.Getenv(tennantID) + "/oauth2/v2.0/token"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, requestURL, bodyParams.Encode())
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

// func main() {
// 	ctx := context.Background()
// 	conf := &oauth2.Config{
// 		ClientID:     os.Getenv(clientID),
// 		ClientSecret: os.Getenv(clientSecret),
// 		Scopes:       []string{"scope", "https://graph.microsoft.com/.default"},
// 		Endpoint: oauth2.Endpoint{
// 			AuthURL:  "https://login.microsoftonline.com/e538d2b6-0142-447b-bf93-985493d12c2e/oauth2/v2.0/authorize",
// 			TokenURL: "https://login.microsoftonline.com/e538d2b6-0142-447b-bf93-985493d12c2e/oauth2/v2.0/token",
// 		},
// 	}

// 	// Redirect user to consent page to ask for permission
// 	// for the scopes specified above.
// 	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
// 	fmt.Printf("Visit the URL for the auth dialog: %v", url)

// 	// Use the authorization code that is pushed to the redirect
// 	// URL. Exchange will do the handshake to retrieve the
// 	// initial access token. The HTTP Client returned by
// 	// conf.Client will refresh the token as necessary.
// 	var code string
// 	if _, err := fmt.Scan(&code); err != nil {
// 		log.Fatal(err)
// 	}
// 	tok, err := conf.Exchange(ctx, code)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	client := conf.Client(ctx, tok)
// 	client.Get("...")
// }
