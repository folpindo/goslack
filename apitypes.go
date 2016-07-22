package goslack

import "fmt"

// APIResponse interface provides common error checking functions
type APIResponse interface {
	Success() bool
	GetError() error
}

// APIType contains the common response fields and behaviour
type APIType struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// Success returns true if no error occured
func (t *APIType) Success() bool {
	return t.OK
}

// GetError returns an error containing details of the fault, if there is one
func (t *APIType) GetError() (err error) {
	if !t.Success() {
		err = fmt.Errorf(t.Error)
	}
	return
}

// RtmStart is the real-time API websocket detail
type RtmStart struct {
	APIType
	URL  string `json:"url"`
	Self Self   `json:"self"`
}

// UserList is a list of users returned from the Slack API
type UserList struct {
	APIType
	Members []UserInfo `json:"members"`
}

// ChannelList is a list of channels returned from the Slack API
type ChannelList struct {
	APIType
	Channels []ChannelInfo `json:"channels"`
}

// GroupList is a list of private channels returned from the Slack API
type GroupList struct {
	APIType
	Groups []ChannelInfo `json:"groups"`
}

// Self is the connection identity
type Self struct {
	ID string `json:"id"`
}

// UserInfo is the detail of a Slack user
type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ChannelInfo is the detail of a Slack channel, group or DM
type ChannelInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsMember bool   `json:"is_member"`
}
