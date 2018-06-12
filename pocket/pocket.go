package pocket

import (
	log "github.com/sirupsen/logrus"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type PocketClient struct {
	key       string
	userToken string
}

func NewClient(key string) *PocketClient {
	return &PocketClient{key: key}
}

func (c *PocketClient) Authenticate() error {
	tokenRequestForm, err := json.Marshal(map[string]string{
		"consumer_key": c.key,
		"redirect_uri": "example.com",
		// TODO: Protect against CSRF by setting a proper value for "state".
		"state": "",
	})
	if err != nil {
		return err
	}

	response, err := http.Post(
		"https://getpocket.com/v3/oauth/request",
		"application/json",
		bytes.NewBuffer(tokenRequestForm),
	)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf(
			"Auth request failed: Pocket responded with %d. Error code %s: %s",
			response.StatusCode,
			response.Header["X-Error-Code"][0],
			response.Header["X-Error"][0],
		)
	}

	rawBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	responseValues, err := url.ParseQuery(string(rawBody))
	if err != nil {
		return err
	}
	response.Body.Close()

	authToken := responseValues.Get("code")

	authUrl := fmt.Sprintf(
		"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
		authToken,
		"example.com",
	)

	fmt.Printf("Visit the URL for the auth dialog: %v\n", authUrl)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return err
	}

	// TODO: Call https://getpocket.com/v3/oauth/authorize

	return nil
}

func (c *PocketClient) Add(url string) (title string, err error) {
	requestBody, err := json.Marshal(map[string]string{
		"url":          url,
		"consumer_key": c.key,
		"access_token": c.userToken,
	})
	if err != nil {
		return "", err
	}

	response, err := http.Post("https://getpocket.com/v3/add",
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return "", err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		log.WithFields(log.Fields{
			"status code":      response.StatusCode,
			"response body":    string(responseBody),
			"response headers": response.Header,
		}).Debug("Pocket returned a non-200 status code")
		return "", fmt.Errorf(
			"Adding an article failed: Pocket responded with %d. Error code %s: %s",
			response.StatusCode,
			response.Header["X-Error-Code"][0],
			response.Header["X-Error"][0],
		)
	}

	responseMap := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &responseMap)
	if err != nil {
		return "", err
	}

	if _, ok := responseMap["item"]; ok {
		title, _ = responseMap["item"].(map[string]interface{})["title"].(string)
	}

	return title, nil
}
