package mangadex

import "encoding/json"

type MangaDexResponse struct {
	// Same as HTTP status code
	Code int `json:"code"`
	// `OK` or `error`
	Status string `json:"status"`
	// Present if Status != OK
	Message string `json:"message"`
	// Actual requested data, only if Status == OK
	Data json.RawMessage `json:"data"`
}

type MangaDexChaptersResponse struct {
	Chapters []MangaChapterList `json:"chapters"`
	Groups   []MangaGroup       `json:"groups"`
}

type MangaRelation struct {
	ID       int    `json:"id"`
	IsHentai bool   `json:"isHentai"`
	Title    string `json:"title"`
	Type     int    `json:"type"`
}

type MangaPublication struct {
	// ???
	Demographic int8 `json:"demographic"`
	// 2 -> Completed
	Status   int8   `json:"status"`
	Language string `json:"language"`
}

type MangaRating struct {
	Bayesian float32 `json:"bayesian"`
	Mean     float32 `json:"mean"`
	Users    int     `json:"users"`
}

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

type MangaCover struct {
	URL    string `json:"url"`
	Volume string `json:"volume"`
}

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

type MangaChapterList struct {
	MangaChapterBase
	Groups []int `json:"groups"`
}

type MangaChapterDetail struct {
	MangaChapterBase
	Groups         []MangaGroup `json:"groups"`
	Status         string       `json:"status"`
	Pages          []string     `json:"pages"`
	Server         string       `json:"server"`
	ServerFallback string       `json:"serverFallback"`
}

type MangaGroupMember struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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
