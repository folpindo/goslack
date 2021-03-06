package goslack

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func getAndUnmarshal(endpoint string, formData url.Values, target APIResponse) (err error) {
	body := formData.Encode()
	url := "https://slack.com" + endpoint + "?" + body
	resp, err := get(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, target)
	if err != nil {
		return
	}

	if !target.Success() {
		err = fmt.Errorf("Error getting response from %s: %s", url, target.GetError())
	}
	return
}

// GetUserList retrieves a list of user IDs mapped to usernames from Slack
func (s *Connection) GetUserList() (users []UserInfo, err error) {
	query := url.Values{}
	query.Add("token", s.Token)

	var userList UserList
	err = getAndUnmarshal("/api/users.list", query, &userList)
	if err != nil {
		return
	}

	users = userList.Members
	return
}

// GetChannelList retrieves a list of active public channels
func (s *Connection) GetChannelList() (channels []ChannelInfo, err error) {
	query := url.Values{}
	query.Add("token", s.Token)
	query.Add("exclude_archived", "1")

	var channelList ChannelList
	err = getAndUnmarshal("/api/channels.list", query, &channelList)
	if err != nil {
		return
	}

	channels = channelList.Channels
	return
}

// GetGroupList retrieves a list of active public channels
func (s *Connection) GetGroupList() (channels []ChannelInfo, err error) {
	query := url.Values{}
	query.Add("token", s.Token)
	query.Add("exclude_archived", "1")

	var groupList GroupList
	err = getAndUnmarshal("/api/groups.list", query, &groupList)
	if err != nil {
		return
	}

	channels = groupList.Groups
	return
}
