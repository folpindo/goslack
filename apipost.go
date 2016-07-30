package goslack

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func postAndUnmarshal(endpoint string, formData url.Values, target APIResponse) (err error) {
	url := "https://slack.com" + endpoint
	resp, err := post(url, formData)
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

// PostIMOpen opens an IM channel with a user
func (s *Connection) PostIMOpen(user string) (channel string, err error) {
	query := url.Values{}
	query.Add("token", s.Token)
	query.Add("user", user)

	response := IMOpenResponse{}
	err = postAndUnmarshal("/api/im.open", query, &response)
	if err != nil {
		return
	}

	return response.Channel.ID, nil
}
