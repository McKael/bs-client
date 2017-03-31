package bsclient

import (
	"errors"
	"net/url"
	"strconv"
)

var (
	errNoSubtitlesFound = errors.New("no subtitles found")
)

// FileName is a string representing a file name in the betaseries API
type FileName string

// Subtitle represents a subtitle returned by the betaseries API
type Subtitle struct {
	ID       int        `json:"id"`
	Language string     `json:"language"`
	Source   string     `json:"source"`
	Quality  int        `json:"quality"`
	File     string     `json:"file"`
	Content  []FileName `json:"content"`
	URL      string     `json:"url"`
	Episode  struct {
		ShowID    int `json:"show_id"`
		EpisodeID int `json:"episode_id"`
		Season    int `json:"season"`
		Episode   int `json:"episode"`
	} `json:"episode"`
	Date string `json:"date"`
}

type subtitles struct {
	Subtitles []Subtitle    `json:"subtitles"`
	Errors    []interface{} `json:"errors"`
}

func (bs *BetaSeries) doGetSubtitles(u *url.URL, usedAPI string) ([]Subtitle, error) {
	resp, err := bs.do("GET", u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := &subtitles{}
	err = bs.decode(data, resp, usedAPI, u.RawQuery)
	if err != nil {
		return nil, err
	}

	if len(data.Subtitles) < 1 {
		return nil, errNoSubtitlesFound
	}

	return data.Subtitles, nil
}

// SubtitlesEpisode returns a slice of subtitles for a given episode
// The language can be provided to filter results (all|vovf|vo|vf).
func (bs *BetaSeries) SubtitlesEpisode(id int, language string) ([]Subtitle, error) {
	usedAPI := "/subtitles/episode"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	if id <= 0 {
		return nil, errIDNotProperlySet
	}
	q.Set("id", strconv.Itoa(id))
	if language != "" {
		q.Set("language", language)
	}
	u.RawQuery = q.Encode()

	return bs.doGetSubtitles(u, usedAPI)
}

// SubtitlesShow returns a slice of subtitles for a given show
// The language can be provided to filter results (all|vovf|vo|vf).
func (bs *BetaSeries) SubtitlesShow(id int, language string) ([]Subtitle, error) {
	usedAPI := "/subtitles/show"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	if id <= 0 {
		return nil, errIDNotProperlySet
	}
	q.Set("id", strconv.Itoa(id))
	if language != "" {
		q.Set("language", language)
	}
	u.RawQuery = q.Encode()

	return bs.doGetSubtitles(u, usedAPI)
}

// SubtitlesLast returns a slice of the last BetaSeries subtitles
// The number can't be higher than 100 with current API.
// The language can be provided to filter results (all|vovf|vo|vf).
func (bs *BetaSeries) SubtitlesLast(number int, language string) ([]Subtitle, error) {
	usedAPI := "/subtitles/last"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	if number > 0 {
		q.Set("number", strconv.Itoa(number))
	}
	if language != "" {
		q.Set("language", language)
	}
	u.RawQuery = q.Encode()

	return bs.doGetSubtitles(u, usedAPI)
}
