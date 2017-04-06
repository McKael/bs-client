package bsclient

import (
	"net/url"
	"strconv"
)

// PlanningGeneral returns a slice of episodes found in [date-before, date+after] timeline.
// Note: the 'date' input must be in YYYY-MM-JJ format or 'now'
// 'eType', the episode type, can be 'premiere' or 'all', or empty.
func (bs *BetaSeries) PlanningGeneral(date, eType string, before, after int) ([]Episode, error) {
	usedAPI := "/planning/general"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	q.Set("date", date)
	q.Set("before", strconv.Itoa(before))
	q.Set("after", strconv.Itoa(after))
	if eType != "" {
		q.Set("type", eType)
	}
	u.RawQuery = q.Encode()
	return bs.doGetEpisodes(u, usedAPI)
}

// PlanningIncoming returns a slice of the first episodes of each tv show
// that are about to be broacasted.
func (bs *BetaSeries) PlanningIncoming() ([]Episode, error) {
	usedAPI := "/planning/incoming"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	return bs.doGetEpisodes(u, usedAPI)
}

// PlanningMember returns a slice of episodes of the member 'id'.
// If 'id' is 0 or negative, the returned episodes are the ones of the
// identified member doing the request.
// The parameter 'unseen' filters not seen episodes.
// The parameter 'month' filters episodes of the given month with the format YYYY-MM.
// Note: the 'month' value can be the string "now".
func (bs *BetaSeries) PlanningMember(id int, unseen bool, month string) ([]Episode, error) {
	usedAPI := "/planning/member"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	if id > 0 {
		q.Set("id", strconv.Itoa(id))
	}
	if unseen {
		q.Set("unseen", "true")
	}
	if month != "" {
		q.Set("month", month)
	}
	u.RawQuery = q.Encode()
	return bs.doGetEpisodes(u, usedAPI)
}
