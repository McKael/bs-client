package bsclient

import (
	"errors"
	"net/url"
	"strconv"
)

var (
	errNoMembersFound = errors.New("no members found")
)

// Member represents the member data returned by the betaserie 'members' API
type Member struct {
	ID     int    `json:"id"`
	FbID   int    `json:"fb_id"`
	Login  string `json:"login"`
	XP     int    `json:"xp"`
	Cached int    `json:"cached"`
	Avatar string `json:"avatar"`
	//ProfileBanner *? `json:"profile_banner"`
	InAccount bool `json:"in_account"`
	Stats     *struct {
		Friends            int     `json:"friends"`
		Shows              int     `json:"shows"`
		Seasons            int     `json:"seasons"`
		Episodes           int     `json:"episodes"`
		Comments           int     `json:"comments"`
		Progress           float64 `json:"progress"`
		EpisodesToWatch    int     `json:"episodes_to_watch"`
		TimeOnTV           int     `json:"time_on_tv"`
		TimeToSpend        int     `json:"time_to_spend"`
		Movies             int     `json:"movies"`
		Badges             int     `json:"badges"`
		MemberSinceDays    int     `json:"member_since_days"`
		FriendsOfFriends   int     `json:"friends_of_friends"`
		EpisodesPerMonth   float64 `json:"episodes_per_month"`
		FavoriteDay        string  `json:"favorite_day"`
		FiveStarsPercent   float64 `json:"five_stars_percent"`
		FourFiveStarsTotal int     `json:"four-five_stars_total"`
		StreakDays         int     `json:"streak_days"`
		FavoriteGenre      string  `json:"favorite_genre"`
		WrittenWords       int     `json:"written_words"`
		WithoutDays        int     `json:"without_days"`
	} `json:"stats"`
	Favorites []Show `json:"favorites"`
	Shows     []Show `json:"shows"`
	Options   *struct {
		Downloaded bool `json:"downloaded"`
		Notation   bool `json:"notation"`
		Timelag    bool `json:"timelag"`
		Global     bool `json:"global"`
		Specials   bool `json:"specials"`
		//EpisodesTri *? `json:"episodes_tri"`
		Friendship string `json:"friendship"`
	} `json:"options"`
}

type members struct {
	Members []Member      `json:"member"`
	Errors  []interface{} `json:"errors"`
}

type memberItem struct {
	Member *Member       `json:"member"`
	Errors []interface{} `json:"errors"`
}

func (bs *BetaSeries) doGetUsers(u *url.URL, usedAPI string) ([]Member, error) {
	resp, err := bs.do("GET", u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users struct {
		Users  []Member      `json:"users"`
		Errors []interface{} `json:"errors"`
	}
	data := &users
	err = bs.decode(data, resp, usedAPI, u.RawQuery)
	if err != nil {
		return nil, err
	}

	if len(data.Users) < 1 {
		return nil, errNoMembersFound
	}

	return data.Users, nil
}

/*
func (bs *BetaSeries) doGetMembers(u *url.URL, usedAPI string) ([]Member, error) {
	resp, err := bs.do("GET", u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := &members{}
	err = bs.decode(data, resp, usedAPI, u.RawQuery)
	if err != nil {
		return nil, err
	}

	if len(data.Members) < 1 {
		return nil, errNoMembersFound
	}

	return data.Members, nil
}
*/

// MembersSearch search for members. 'login' can contain the wildcard '%'
func (bs *BetaSeries) MembersSearch(login string, limit int) ([]Member, error) {
	usedAPI := "/members/search"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	q.Set("login", login)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	u.RawQuery = q.Encode()

	return bs.doGetUsers(u, usedAPI)
}

// MembersInfos returns member information about the given user
// If summary is true, no data about movies and shows is returns.
// If summary is false, only can optionally be set to 'movies' or 'shows'.
func (bs *BetaSeries) MembersInfos(id int, summary bool, only string) (*Member, error) {
	usedAPI := "/members/infos"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	if id < 1 {
		return nil, errIDNotProperlySet
	}
	q.Set("id", strconv.Itoa(id))
	if summary {
		q.Set("summary", "true")
	} else if only != "" {
		q.Set("only", only)
	}
	u.RawQuery = q.Encode()

	resp, err := bs.do("GET", u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := &memberItem{}
	err = bs.decode(data, resp, usedAPI, u.RawQuery)
	if err != nil {
		return nil, err
	}

	return data.Member, nil
}
