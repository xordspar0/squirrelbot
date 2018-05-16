package pocket

import (
	log "github.com/sirupsen/logrus"

	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type PocketClient struct {
	key       string
	userToken string
}

func NewClient(key, userToken string) *PocketClient {
	return &PocketClient{
		key:       key,
		userToken: userToken,
	}
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
			"status code": response.StatusCode,
			"response body": responseBody,
		}).Debug("Pocket returned a non-200 status code")

		return "", errors.New("Pocket returned " + string(responseBody))
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
