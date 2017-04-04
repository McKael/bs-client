package bsclient

import (
	"net/url"
	"strconv"
)

func (bs *BetaSeries) friendUpdate(method, endpoint string, id int) (*Member, error) {
	usedAPI := "/friends/" + endpoint
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

	friend := &memberItem{}
	err = bs.decode(friend, resp, usedAPI, u.RawQuery)
	if err != nil {
		return nil, err
	}

	return friend.Member, nil
}

// FriendsList lists a member's friends
// If 'blocked' is true, return the list of blocked users (only if id not set)
func (bs *BetaSeries) FriendsList(id int, blocked bool) ([]Member, error) {
	usedAPI := "/friends/list"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	if id > 0 {
		q.Set("id", strconv.Itoa(id))
	}
	if blocked {
		q.Set("blocked", "true")
	}
	u.RawQuery = q.Encode()

	return bs.doGetUsers(u, usedAPI)
}

// FriendsRequests returns a list of members the user has sent friendship requests to
// If 'received' is true, returns a list of members that sent friendship requests
func (bs *BetaSeries) FriendsRequests(received bool) ([]Member, error) {
	usedAPI := "/friends/requests"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	if received {
		q.Set("received", "true")
	}
	u.RawQuery = q.Encode()

	return bs.doGetUsers(u, usedAPI)
}

// FriendsFriend adds the member 'id' to the user account
func (bs *BetaSeries) FriendsFriend(id int) (*Member, error) {
	return bs.friendUpdate("POST", "friend", id)
}

// FriendsNotFriend removes the member 'id' from the user account
func (bs *BetaSeries) FriendsNotFriend(id int) (*Member, error) {
	return bs.friendUpdate("DELETE", "friend", id)
}

// FriendsBlock blocks the member 'id'
func (bs *BetaSeries) FriendsBlock(id int) (*Member, error) {
	return bs.friendUpdate("POST", "block", id)
}

// FriendsUnblock unblocks the member 'id'
func (bs *BetaSeries) FriendsUnblock(id int) (*Member, error) {
	return bs.friendUpdate("DELETE", "block", id)
}
