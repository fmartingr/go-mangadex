package mangadex

import (
	"encoding/json"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

// GetChapters - Requests the chapters and groups for the provided manga instance.
// Requires a ChaptersParams argument to work through the pagination and other request
// parameters.
func (manga *Manga) GetChapters(params ChaptersParams) ([]MangaChapterList, []MangaGroup, error) {
	var mangaChaptersResult []MangaChapterList
	var mangaGroupsResult []MangaGroup
	params.validate()

	response, errRequest := doRequest("GET", APIBaseURL+path.Join("manga", strconv.Itoa(manga.ID), "chapters")+"?"+params.asQueryParams().Encode())
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

// GetChapter retrieves the specific chapter detail from the provided Manga.
// This function returns a more detailed chapter object since the list only returns information
// but the detail endpoint is needed in order to get the pages and servers where those are stored.
func (manga *Manga) GetChapter(chapter string) (MangaChapterDetail, error) {
	var result MangaChapterDetail

	response, errRequest := doRequest("GET", APIBaseURL+path.Join("chapter", chapter))
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

// GetCovers requests the covers for the provided manga
func (manga *Manga) GetCovers() ([]MangaCover, error) {
	var result []MangaCover
	response, errRequest := doRequest("GET", APIBaseURL+path.Join("manga", strconv.Itoa(manga.ID), "covers"))
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

// GetManga retrieves the manga information for the provided ID.
func GetManga(mangaID int) (Manga, error) {
	result := Manga{}
	response, errRequest := doRequest("GET", APIBaseURL+path.Join("manga", strconv.Itoa(mangaID)))
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
