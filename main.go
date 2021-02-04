package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const APIBaseURL = "https://api.mangadex.org/v2/"

type MangaDexResponse struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

type MangaRelation struct {
	ID       int    `json:"id"`
	IsHentai bool   `json:"isHentai"`
	Title    string `json:"title"`
	Type     int    `json:"type"`
}

type Manga struct {
	ID                int      `json:"id"`
	AlternativeTitles []string `json:"altTitles"`
	Artist            []string `json:"artist"`
	Author            []string `json:"author"`
	Comments          int      `json:"comments"`
	Description       string   `json:"description"`
	Follows           int      `json:"follows"`
	IsHentai          bool     `json:"isHentai"`
	LastChapter       string   `json:"lastChapter"`
	LastUploaded      int64    `json:"lastUploaded"`
	LastVolume        string   `json:"lastVolume"`
	// Links // key->value, fixed keys
	Cover string `json:"mainCover"`
	// Publication
	// Rating
	// Relations
	// Tags
	Title string `json:"title"`
	Views int    `json:"views"`
}

type MangaCover struct {
	URL    string `json:"url"`
	Volume string `json:"volume"`
}

func DoRequest(method string, requestURL string) (*MangaDexResponse, error) {
	result := MangaDexResponse{}
	parsedURL, errParse := url.Parse(requestURL)
	if errParse != nil {
		return &result, errParse
	}

	log.Printf("[debug] request: %s", parsedURL)
	request := http.Request{
		Method: method,
		URL:    parsedURL,
		Header: map[string][]string{
			"User-Agent": {"go-mangadex/0.0.1"},
		},
	}

	response, errResponse := http.DefaultClient.Do(&request)
	if errResponse != nil {
		log.Print(errResponse)
		return &result, errResponse
	}

	if response.StatusCode != 200 {
		log.Printf("[error] Status code != 200 -- %d", response.StatusCode)
		return &result, errors.New(strconv.Itoa(response.StatusCode))
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	body, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("[error] %s", errRead)
		return &result, errRead
	}

	errJSON := json.Unmarshal(body, &result)
	if errJSON != nil {
		log.Print(errJSON)
		return &result, errJSON
	}

	return &result, nil
}

func (manga *Manga) GetCovers() ([]MangaCover, error) {
	var result []MangaCover
	response, errRequest := DoRequest("GET", APIBaseURL+path.Join("manga", strconv.Itoa(manga.ID), "covers"))
	if errRequest != nil {
		return result, errRequest
	}

	errJSON := json.Unmarshal(response.Data, &result)
	if errJSON != nil {
		log.Printf("[error] %s", errJSON)
	}
	return result, nil
}

func GetManga(mangaID int) (Manga, error) {
	result := Manga{}
	response, errRequest := DoRequest("GET", APIBaseURL+path.Join("manga", strconv.Itoa(mangaID)))
	if errRequest != nil {
		return result, errRequest
	}
	//log.Print(response)

	errJSON := json.Unmarshal(response.Data, &result)
	if errJSON != nil {
		log.Printf("[error] %s", errJSON)
	}
	return result, nil
}

// func (manga *Manga)

func main() {
	manga, err := GetManga(2890)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", manga)
	covers, err := manga.GetCovers()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", covers)
}
