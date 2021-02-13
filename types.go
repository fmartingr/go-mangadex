package mangadex

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Response handles the response from MangaDex
type Response struct {
	// Same as HTTP status code
	Code int `json:"code"`
	// `OK` or `error`
	Status string `json:"status"`
	// Present if Status != OK
	Message string `json:"message"`
	// Actual requested data, only if Status == OK
	Data json.RawMessage `json:"data"`
}

// IsOK Checks if the response is correct
func (response Response) IsOK() bool {
	return response.Status == "OK"
}

// ChaptersResponse handles the response of the chapters list which returns
// two kinds of objects in the `json:"data"` key.
type ChaptersResponse struct {
	Chapters []MangaChapterList `json:"chapters"`
	Groups   []MangaGroup       `json:"groups"`
}

// MangaRelation relations between mangas
type MangaRelation struct {
	ID       int    `json:"id"`
	IsHentai bool   `json:"isHentai"`
	Title    string `json:"title"`
	Type     int    `json:"type"`
}

// MangaPublication stores certain information for the publication of the manga
// Some values are easily guessed, others are not ...
type MangaPublication struct {
	// ???
	Demographic int8 `json:"demographic"`
	// 2 -> Completed
	Status   int8   `json:"status"`
	Language string `json:"language"`
}

// IsComplete returns if the manga has finished publishing
func (publication MangaPublication) IsComplete() bool {
	return publication.Status == 2
}

// MangaRating the rating for a particular manga
type MangaRating struct {
	Bayesian float32 `json:"bayesian"`
	Mean     float32 `json:"mean"`
	Users    int     `json:"users"`
}

// MangaLinks contains the relation of the links map with more verbose names
type MangaLinks struct {
	AniList      string `json:"al"`
	AnimePlanet  string `json:"ap"`
	BookWalker   string `json:"bw"`
	Kitsu        string `json:"kt"`
	MangaUpdates string `json:"mu"`
	Amazon       string `json:"amz"`
	EBookJapan   string `json:"ebj"`
	MyAnimeList  string `json:"mal"`
	Raw          string `json:"raw"`
	EnglishRaw   string `json:"engtl"`
}

// Manga stores the base manga object and details
type Manga struct {
	ID                int              `json:"id"`
	AlternativeTitles []string         `json:"altTitles"`
	Artist            []string         `json:"artist"`
	Author            []string         `json:"author"`
	Comments          int              `json:"comments"`
	Description       string           `json:"description"`
	Follows           int              `json:"follows"`
	IsHentai          bool             `json:"isHentai"`
	Links             MangaLinks       `json:"links"`
	LastChapter       string           `json:"lastChapter"`
	LastUploaded      int64            `json:"lastUploaded"`
	LastVolume        string           `json:"lastVolume"`
	Cover             string           `json:"mainCover"`
	Publication       MangaPublication `json:"publication"`
	Rating            MangaRating      `json:"rating"`
	Relations         []MangaRelation  `json:"relations"`
	Tags              []int16          `json:"tags"`
	Title             string           `json:"title"`
	Views             int              `json:"views"`
}

// MangaCover stores the cover object
type MangaCover struct {
	URL    string `json:"url"`
	Volume string `json:"volume"`
}

// MangaChapterBase the base attributes for a chapter for both the list and the detail
type MangaChapterBase struct {
	ID         int    `json:"id"`
	Hash       string `json:"hash"`
	MangaID    int    `json:"mangaId"`
	MangaTitle string `json:"mangaTitle"`
	Volume     string `json:"volume"`
	Chapter    string `json:"chapter"`
	Title      string `json:"title"`
	Language   string `json:"language"`
	Uploader   int    `json:"uploader"`
	Timestamp  int64  `json:"timestamp"`
	ThreadID   int    `json:"threadId"`
	Comments   int    `json:"comments"`
	Views      int    `json:"views"`
}

// MangaChapterList stores the chapter object from the listing, which only uses the
// base details and uses the ID for the groups instead of returning the entire object
type MangaChapterList struct {
	MangaChapterBase
	Groups []int `json:"groups"`
}

// MangaChapterDetail stores the bases of a chapter plus the required attributes
// to retrieve the page blobs
type MangaChapterDetail struct {
	MangaChapterBase
	Groups         []MangaGroup `json:"groups"`
	Status         string       `json:"status"`
	Pages          []string     `json:"pages"`
	Server         string       `json:"server"`
	ServerFallback string       `json:"serverFallback"`
}

// MangaGroupMember stores the member of a group
type MangaGroupMember struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MangaGroup stores the group behind releases of a particular manga
type MangaGroup struct {
	ID               int                `json:"id"`
	Name             string             `json:"name"`
	AlternativeNames string             `json:"altNames"`
	Language         string             `json:"language"`
	Leader           MangaGroupMember   `json:"leader"`
	Members          []MangaGroupMember `json:"members"`
	Description      string             `json:"description"`
	Website          string             `json:"website"`
	Discord          string             `json:"discord"`
	IRCServer        string             `json:"ircServer"`
	IRCChannel       string             `json:"ircChannel"`
	EMail            string             `json:"email"`
	Founded          string             `json:"founded"`
	Likes            int                `json:"likes"`
	Follows          int                `json:"follows"`
	Views            int                `json:"views"`
	Chapters         int                `json:"chapters"`
	ThreadID         int                `json:"threadId"`
	ThreadPosts      int                `json:"threadPosts"`
	IsLocked         bool               `json:"isLocked"`
	IsInactive       bool               `json:"isInactive"`
	Delay            int                `json:"delay"`
	LastUpdated      int                `json:"lastUpdated"`
	Banner           string             `json:"banner"`
}

// ChaptersParams the request parameters for the chapters listing
type ChaptersParams struct {
	// How many items per page (max 100)
	Limit int `json:"limit"`
	// Page to retrieve (default 0)
	Page int `json:"p"`
	// Hide groups blocked by the user (auth not implemented)
	BlockGroups bool `json:"blockgroups"`
}

// NewChaptersParams returns a ChapterParams object with sensible defaults
func NewChaptersParams() ChaptersParams {
	return ChaptersParams{
		Limit:       100,
		Page:        0,
		BlockGroups: false,
	}
}

func (params *ChaptersParams) validate() {
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 100
	}
}

func (params *ChaptersParams) asQueryParams() url.Values {
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
