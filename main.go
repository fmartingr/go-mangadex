package mangadex

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

const APIBaseURL = "https://api.mangadex.org/v2/"

var cacheEnabled bool = false

func EnableCache() {
	cacheEnabled = true
}

func DisableCache() {
	cacheEnabled = false
}

func getCachePath() string {
	userCacheDir, errCache := os.UserCacheDir()
	if errCache != nil {
		logrus.Fatalf("Unable to retrieve cache directory: %s", errCache)
	}
	return filepath.Join(userCacheDir, "go-mangadex")
}

func getCachePathFor(mangadexURL string) string {
	fileName := getCacheFilename(mangadexURL)
	return filepath.Join(getCachePath(), fileName)
}

func getCacheFilename(mangadexURL string) string {
	urlHash := sha1.New()
	_, errWrite := urlHash.Write([]byte(mangadexURL))
	if errWrite != nil {
		logrus.Errorf("Error generating hash for %s: %s", mangadexURL, errWrite)
	}
	return fmt.Sprintf("%x", urlHash.Sum(nil))
}

func cacheExists(mangadexURL string) bool {
	stat, err := os.Stat(getCachePathFor(mangadexURL))
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}

func initCache() {
	cachePath := getCachePath()
	_, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		logrus.Infof("Cache directory does not exist, creating. [%s]", cachePath)
		errCachePath := os.MkdirAll(cachePath, 0755)
		if errCachePath != nil {
			logrus.Errorf("Cache directory couldn't be generated, caching will likely fail: %s", errCachePath)
		}
	}
}

func DoRequest(method string, requestURL string) (*MangaDexResponse, error) {
	result := MangaDexResponse{}
	parsedURL, errParse := url.Parse(requestURL)
	if errParse != nil {
		return &result, errParse
	}

	if cacheEnabled {
		initCache()
		if cacheExists(parsedURL.String()) {
			cacheData, errRead := ioutil.ReadFile(getCachePathFor(parsedURL.String()))
			if errRead != nil {
				logrus.Fatalf("Error reading cache for URL: %s [%s]: %s", parsedURL.String(), getCacheFilename(parsedURL.String()), errRead)
			}
			errJSON := json.Unmarshal(cacheData, &result)
			if errJSON != nil {
				logrus.Fatalf("Error parsing JSON from cache: %s: %s", getCacheFilename(parsedURL.String()), errJSON)
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
	if cacheEnabled {
		logrus.Infof("Writting cache for %s", parsedURL.String())
		logrus.Infof("Writting cache to: %s", getCacheFilename(parsedURL.String()))
		errWriteCache := ioutil.WriteFile(getCachePathFor(parsedURL.String()), body, 0644)
		if errWriteCache != nil {
			logrus.Warnf("Can't write to cache: %s", errWriteCache)
		}
	}

	errJSON := json.Unmarshal(body, &result)
	if errJSON != nil {
		logrus.Errorf("Error parsing body: %s", errJSON)
		return &result, errJSON
	}

	return &result, nil
}

func (manga *Manga) GetCovers() ([]MangaCover, error) {
	var result []MangaCover
	response, errRequest := DoRequest("GET", APIBaseURL+path.Join("manga", strconv.Itoa(manga.ID), "covers"))
	if errRequest != nil {
		logrus.Errorf("Request error: %s", errRequest)
		return result, errRequest
	}

	errJSON := json.Unmarshal(response.Data, &result)
	if errJSON != nil {
		logrus.Errorf("Error parsing JSON: %s", errJSON)
		return result, errJSON
	}
	return result, nil
}

type GetChaptersParams struct {
	Limit       int  `json:"limit"`
	Page        int  `json:"p"`
	BlockGroups bool `json:"blockgroups"`
}

func NewGetChaptersParams() GetChaptersParams {
	return GetChaptersParams{
		Limit:       100,
		Page:        0,
		BlockGroups: false,
	}
}

func (params *GetChaptersParams) Validate() {
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 100
	}
}

func (params *GetChaptersParams) AsQueryParams() url.Values {
	queryParams := url.Values{}

	if params.Page > 0 {
		queryParams.Add("p", strconv.FormatInt(int64(params.Page), 10))
	}
	queryParams.Add("limit", strconv.FormatInt(int64(params.Limit), 10))
	if params.BlockGroups {
		queryParams.Add("blockgroups", strconv.FormatBool(params.BlockGroups))
	}

	return queryParams
}

func (manga *Manga) GetChapters(params GetChaptersParams) ([]MangaChapter, []MangaGroup, error) {
	var mangaChaptersResult []MangaChapter
	var mangaGroupsResult []MangaGroup
	params.Validate()

	response, errRequest := DoRequest("GET", APIBaseURL+path.Join("manga", strconv.Itoa(manga.ID), "chapters")+"?"+params.AsQueryParams().Encode())
	if errRequest != nil {
		logrus.Errorf("Request error: %s", errRequest)
		return mangaChaptersResult, mangaGroupsResult, errRequest
	}

	var mangaDexChaptersResponse MangaDexChaptersResponse

	errJSON := json.Unmarshal(response.Data, &mangaDexChaptersResponse)
	if errJSON != nil {
		logrus.Errorf("Error parsing JSON: %s", errJSON)
		return mangaChaptersResult, mangaGroupsResult, errJSON
	}

	return mangaDexChaptersResponse.Chapters, mangaDexChaptersResponse.Groups, nil
}

func (manga *Manga) GetChapter(chapter string) (MangaChapter, error) {
	var result MangaChapter

	response, errRequest := DoRequest("GET", APIBaseURL+path.Join("chapter", chapter))
	if errRequest != nil {
		logrus.Errorf("Request error: %s", errRequest)
		return result, errRequest
	}

	errJSON := json.Unmarshal(response.Data, &result)
	if errJSON != nil {
		logrus.Errorf("Error parsing JSON: %s", errJSON)
		return result, errJSON
	}

	return result, nil
}

func (manga *Manga) GetVolumeChapters(volume string) ([]MangaChapter, error) {
	var result []MangaChapter

	return result, nil
}

func GetManga(mangaID int) (Manga, error) {
	result := Manga{}
	response, errRequest := DoRequest("GET", APIBaseURL+path.Join("manga", strconv.Itoa(mangaID)))
	if errRequest != nil {
		logrus.Errorf("Request error: %s", errRequest)
		return result, errRequest
	}

	errJSON := json.Unmarshal(response.Data, &result)
	if errJSON != nil {
		logrus.Errorf("Error parsing JSON: %s", errJSON)
		return result, errJSON
	}
	return result, nil
}
