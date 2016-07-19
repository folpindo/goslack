package rtmtypes

import (
	"encoding/json"
	"fmt"

	"github.com/doozr/goslack/apitypes"
)

// RtmEvent is the base Slack event that others extend
type RtmEvent struct {
	Type string `json:"type"`
}

// RtmRaw is the raw JSON data and type tag that can be coerced
type RtmRaw struct {
	RtmEvent
	Raw []byte
}

// RtmMessage is a Slack message sent to a channel
type RtmMessage struct {
	RtmEvent
	ID      uint64 `json:"id"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
}

// RtmUserChange is a notification that a user profile has changed
type RtmUserChange struct {
	RtmEvent
	User apitypes.UserInfo `json:"user"`
}

// Unmarshal creates an RtmEvent instance from raw bytes
func Unmarshal(raw []byte) (event RtmRaw, err error) {
	err = json.Unmarshal(raw, &event)
	if err != nil {
		return
	}

	if event.Type == "" {
		err = fmt.Errorf("Missing type in Slack event `%v`", raw)
	} else {
		event.Raw = raw
	}
	return
}

func (e *RtmRaw) unmarshalToType(t string, v interface{}) (err error) {
	if e.Type != t {
		err = fmt.Errorf("Expected type `%s` but received `%s`", t, e.Type)
	} else {
		err = json.Unmarshal(e.Raw, v)
	}
	return
}

// RtmMessage unmarshals an RtmRaw into an RtmMessage
func (e *RtmRaw) RtmMessage() (event RtmMessage, err error) {
	err = e.unmarshalToType("message", &event)
	return
}

// RtmUserChange unmarshals an RtmRaw into an RtmUserChange
func (e *RtmRaw) RtmUserChange() (event RtmUserChange, err error) {
	err = e.unmarshalToType("user_change", &event)
	return
}
