package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/doozr/goslack/apitypes"
	"github.com/doozr/goslack/rtmtypes"

	"golang.org/x/net/websocket"
)

// RtmChannel is an incoming stream of messages on the Slack websocket
type RtmChannel chan rtmtypes.RtmRaw

// Slack represents a connection to the Slack web and real time APIs
type Slack struct {
	Token     string
	RealTime  RtmChannel
	ID        string
	websocket *websocket.Conn
}

// Error is any error raised due to a Slack fault or API problem
type Error struct {
	msg string
}

// New creates a new Slack instance
func New(token string) (slackConn *Slack, err error) {
	wsurl, id, err := getWebsocketURL(token)
	if err != nil {
		return
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		return
	}

	slackConn = &Slack{
		Token:     token,
		RealTime:  make(RtmChannel, 10),
		ID:        id,
		websocket: ws,
	}

	go getRealTimeEvents(slackConn)

	return
}

func getRealTimeEvents(s *Slack) {
	for {
		var data []byte
		err := websocket.Message.Receive(s.websocket, &data)
		if err != nil {
			log.Printf("Error receiving from websocket: %q", err)
			continue
		}

		e, err := rtmtypes.Unmarshal(data)
		if err != nil {
			log.Printf("Error unmarshalling message: %q", err)
			continue
		}

		s.RealTime <- e
	}
}

var counter uint64

// PostMessage sends a message to a Slack channel
func (s *Slack) PostMessage(channel, text string) error {
	id := atomic.AddUint64(&counter, 1)
	m := rtmtypes.RtmMessage{
		RtmEvent: rtmtypes.RtmEvent{Type: "message"},
		ID:       id,
		Channel:  channel,
		User:     "",
		Text:     text,
	}
	return websocket.JSON.Send(s.websocket, m)
}

// GetUserList retrieves a list of user IDs mapped to usernames from Slack
func (s *Slack) GetUserList() (users []apitypes.UserInfo, err error) {
	body := encodeFormData(map[string]string{
		"token": s.Token,
	})

	resp, err := get("https://slack.com/api/users.list?" + body)
	if err != nil {
		return
	}

	var response apitypes.UserList
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return
	}

	if !response.Ok {
		err = fmt.Errorf("Error getting user info: %s", response.Error)
	}

	users = response.Members
	return
}
