package mangadex

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

const APIBaseURL = "https://api.mangadex.org/v2/"

func doRequest(method string, requestURL string) (*MangaDexResponse, error) {
	result := MangaDexResponse{}
	parsedURL, errParse := url.Parse(requestURL)
	if errParse != nil {
		return &result, errParse
	}

	if cacheEnabled {
		initCache()
		if cacheExists(parsedURL) {
			cacheData, errRead := ioutil.ReadFile(getCachePathFor(parsedURL))
			if errRead != nil {
				logrus.Fatalf("Error reading cache for URL: %s [%s]: %s", parsedURL.String(), getCacheFilename(parsedURL), errRead)
			}
			errJSON := json.Unmarshal(cacheData, &result)
			if errJSON != nil {
				logrus.Fatalf("Error parsing JSON from cache: %s: %s", getCacheFilename(parsedURL), errJSON)
			}
			logrus.Debugf("Request loaded from cache: %s", parsedURL.String())
			return &result, nil
		} else {
			logrus.Debugf("Cache not found for %s", parsedURL.String())
		}
	}

	logrus.Tracef("Making request %s", parsedURL)
	request := http.Request{
		Method: method,
		URL:    parsedURL,
		Header: map[string][]string{
			"User-Agent": {"go-mangadex/0.0.1"},
		},
	}

	response, errResponse := http.DefaultClient.Do(&request)
	if errResponse != nil {
		logrus.Tracef("Request error: %s", errResponse)
		return &result, errResponse
	}

	if response.StatusCode != 200 {
		logrus.Tracef("Response status code not successful: %d", response.StatusCode)
		logrus.Tracef("Response body: %s", response.Body)
		return &result, errors.New(strconv.Itoa(response.StatusCode))
	}

	logrus.Tracef("Response status code: %s", response.Status)

	if response.Body != nil {
		defer response.Body.Close()
	}

	body, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		logrus.Errorf("Error reading body: %s", errRead)
		return &result, errRead
	}

	logrus.Tracef("Response body: %s", body)

	// Write cache
	logrus.Infof("Writting cache for %s", parsedURL.String())
	logrus.Infof("Writting cache to: %s", getCacheFilename(parsedURL))
	errWriteCache := ioutil.WriteFile(getCachePathFor(parsedURL), body, 0644)
	if errWriteCache != nil {
		logrus.Warnf("Can't write to cache: %s", errWriteCache)
	}

	errJSON := json.Unmarshal(body, &result)
	if errJSON != nil {
		logrus.Errorf("Error parsing body: %s", errJSON)
		return &result, errJSON
	}

	return &result, nil
}
