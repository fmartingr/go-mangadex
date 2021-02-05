package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const APIBaseURL = "https://api.mangadex.org/v2/"

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

func (manga *Manga) GetChapters(params GetChaptersParams) ([]MangaChapter, error) {
	var result []MangaChapter
	params.Validate()

	response, errRequest := DoRequest("GET", APIBaseURL+path.Join("manga", strconv.Itoa(manga.ID), "chapters")+"?"+params.AsQueryParams().Encode())
	if errRequest != nil {
		return result, errRequest
	}

	var mangaDexChaptersResponse MangaDexChaptersResponse

	errJSON := json.Unmarshal(response.Data, &mangaDexChaptersResponse)
	if errJSON != nil {
		return result, errJSON
	}
	result = mangaDexChaptersResponse.Chapters

	return result, nil
}

func (manga *Manga) GetChapter(chapter string) (MangaChapter, error) {
	var result MangaChapter

	response, errRequest := DoRequest("GET", APIBaseURL+path.Join("chapter", chapter))
	if errRequest != nil {
		return result, errRequest
	}

	errJSON := json.Unmarshal(response.Data, &result)
	if errJSON != nil {
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
		return result, errRequest
	}

	errJSON := json.Unmarshal(response.Data, &result)
	if errJSON != nil {
		log.Printf("[error] %s", errJSON)
	}
	return result, nil
}

func main() {
	manga, err := GetManga(2890)
	if err != nil {
		panic(err)
	}

	log.Printf("Manga: %s", manga.Title)

	chapters, err := manga.GetChapters(NewGetChaptersParams())
	if err != nil {
		panic(err)
	}

	for i := range chapters {
		log.Printf("v%2s %3s \t %s", chapters[i].Volume, chapters[i].Chapter, chapters[i].Title)
	}
}
