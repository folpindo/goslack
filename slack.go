package goslack

import (
	"log"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

// RtmChannel is an incoming stream of messages on the Slack websocket
type RtmChannel chan RtmRaw

// Connection represents a connection to the Slack web and real time APIs
type Connection struct {
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
func New(token string) (connection *Connection, err error) {
	wsurl, id, err := getWebsocketURL(token)
	if err != nil {
		return
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		return
	}

	connection = &Connection{
		Token:     token,
		RealTime:  make(RtmChannel, 10),
		ID:        id,
		websocket: ws,
	}

	go getRealTimeEvents(connection)

	return
}

func getRealTimeEvents(s *Connection) {
	for {
		var data []byte
		err := websocket.Message.Receive(s.websocket, &data)
		if err != nil {
			log.Printf("Error receiving from websocket: %q", err)
			continue
		}

		e, err := unmarshal(data)
		if err != nil {
			continue
		}

		s.RealTime <- e
	}
}

var counter uint64

// PostRealTimeMessage sends a message to a Slack channel
func (s *Connection) PostRealTimeMessage(channel, text string) error {
	id := atomic.AddUint64(&counter, 1)
	m := RtmMessage{
		RtmEvent: RtmEvent{Type: "message"},
		ID:       id,
		Channel:  channel,
		User:     "",
		Text:     text,
	}
	return websocket.JSON.Send(s.websocket, m)
}
