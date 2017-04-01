package bsclient

import (
	"net/url"
	"strconv"
)

type episodeItem struct {
	Episode *Episode      `json:"episode"`
	Errors  []interface{} `json:"errors"`
}

// episodeGet returns an episode
// Note: scraper and list cannot be requested with this method
func (bs *BetaSeries) episodeGet(endPoint string, id, theTvdbID int,
	subtitles bool, number string) (*Episode, error) {
	// endPoint can be: display, latest, next, search
	usedAPI := "/episodes/" + endPoint
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()

	if endPoint == "search" {
		if id > 0 {
			q.Set("show_id", strconv.Itoa(id))
		}
		if number != "" {
			q.Set("number", number)
		}
	} else if id > 0 {
		q.Set("id", strconv.Itoa(id))
	}

	if theTvdbID > 0 {
		q.Set("thetvdb_id", strconv.Itoa(theTvdbID))
	}

	if subtitles {
		q.Set("subtitles", "true")
	}

	u.RawQuery = q.Encode()

	resp, err := bs.do("GET", u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	episode := &episodeItem{}
	err = bs.decode(episode, resp, usedAPI, u.RawQuery)
	if err != nil {
		return nil, err
	}

	return episode.Episode, nil
}

func (bs *BetaSeries) episodeUpdate(method, endpoint string, id, theTvdbID int) (*Episode, error) {
	usedAPI := "/episodes/" + endpoint
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	q.Set("id", strconv.Itoa(id))
	u.RawQuery = q.Encode()

	resp, err := bs.do(method, u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	episode := &episodeItem{}
	err = bs.decode(episode, resp, usedAPI, u.RawQuery)
	if err != nil {
		return nil, err
	}

	return episode.Episode, nil
}

// EpisodeScraper returns an episode from a file name
func (bs *BetaSeries) EpisodeScraper(fileName string) (*Episode, error) {
	usedAPI := "/episodes/scraper"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	q.Set("file", fileName)
	u.RawQuery = q.Encode()

	resp, err := bs.do("GET", u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	episode := &episodeItem{}
	err = bs.decode(episode, resp, usedAPI, u.RawQuery)
	if err != nil {
		return nil, err
	}

	return episode.Episode, nil
}

// EpisodeLatest returns the latest episode for a given show
func (bs *BetaSeries) EpisodeLatest(showID, theTvdbShowID int) (*Episode, error) {
	return bs.episodeGet("latest", showID, theTvdbShowID, false, "")
}

// EpisodeDisplay returns the latest episode for a given show
func (bs *BetaSeries) EpisodeDisplay(showID, theTvdbShowID int, subtitles bool) (*Episode, error) {
	return bs.episodeGet("display", showID, theTvdbShowID, subtitles, "")
}

// EpisodeNext returns the next episode for a given show
func (bs *BetaSeries) EpisodeNext(showID, theTvdbShowID int) (*Episode, error) {
	return bs.episodeGet("next", showID, theTvdbShowID, false, "")
}

// EpisodeSearch returns an episode for a given show based on its number
func (bs *BetaSeries) EpisodeSearch(showID int, subtitles bool, number string) (*Episode, error) {
	return bs.episodeGet("search", showID, 0, subtitles, number)
}

// EpisodeDownloaded marks the episode with the given id as downloaded.
func (bs *BetaSeries) EpisodeDownloaded(bsID, theTvdbID int) (*Episode, error) {
	return bs.episodeUpdate("POST", "downloaded", bsID, theTvdbID)
}

// EpisodeNotDownloaded marks the episode with the given id as not downloaded.
func (bs *BetaSeries) EpisodeNotDownloaded(bsID, theTvdbID int) (*Episode, error) {
	return bs.episodeUpdate("DELETE", "downloaded", bsID, theTvdbID)
}

// EpisodeWatched marks the episode with the given id as watched.
func (bs *BetaSeries) EpisodeWatched(bsID, theTvdbID int) (*Episode, error) {
	return bs.episodeUpdate("POST", "watched", bsID, theTvdbID)
}

// EpisodeNotWatched marks the episode with the given id as not watched.
func (bs *BetaSeries) EpisodeNotWatched(bsID, theTvdbID int) (*Episode, error) {
	return bs.episodeUpdate("DELETE", "watched", bsID, theTvdbID)
}
