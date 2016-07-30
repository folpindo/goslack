package goslack

import (
	"log"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

// PingTimeout is the time between ping requests
const PingTimeout = 10 * time.Second

// ActivityTimeout is the idle time before a reconnect is attempted
const ActivityTimeout = 30 * time.Second

// RtmChannel is an incoming stream of messages on the Slack websocket
type RtmChannel chan RtmRaw

// Connection represents a connection to the Slack web and real time APIs
type Connection struct {
	Token     string
	RealTime  RtmChannel
	ID        string
	websocket *websocket.Conn
}

// New creates a new Slack instance
func New(token string) (connection *Connection, err error) {
	ws, id, err := connectWebsocket(token)
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

func connectWebsocket(token string) (ws *websocket.Conn, id string, err error) {
	wsurl, id, err := getWebsocketURL(token)
	if err != nil {
		return
	}

	ws, err = websocket.Dial(wsurl, "", "https://api.slack.com/")
	return
}

func reconnectWebsocket(connection *Connection) (err error) {
	ws, id, err := connectWebsocket(connection.Token)
	if err != nil {
		return
	}

	connection.ID = id
	connection.websocket = ws
	return
}

func getRealTimeEvents(s *Connection) {
	lastActivity := time.Now()
	rtmEvents := make(chan RtmRaw)
	go getWebsocketEvent(s, rtmEvents)

	for {
		select {
		case e := <-rtmEvents:
			lastActivity = time.Now()
			s.RealTime <- e
		case <-time.After(PingTimeout):
			if time.Since(lastActivity) >= ActivityTimeout {
				log.Printf("Activity timeout after %s: reconnecting", ActivityTimeout)
				err := reconnectWebsocket(s)
				if err != nil {
					log.Printf("Error while reconnecting: %s", err)
				}
			} else {
				s.Ping()
			}
		}
	}
}

func getWebsocketEvent(s *Connection, ch chan RtmRaw) {
	for {
		var data []byte
		err := websocket.Message.Receive(s.websocket, &data)
		if err != nil {
			log.Printf("Error receiving from websocket: %q", err)
		}

		e, err := unmarshal(data)
		if err == nil {
			ch <- e
		}
	}
}

var counter uint64

func nextID() uint64 {
	return atomic.AddUint64(&counter, 1)
}

// PostRealTimeMessage sends a message to a Slack channel
func (s *Connection) PostRealTimeMessage(channel, text string) error {
	id := nextID()
	m := RtmMessage{
		RtmEvent: RtmEvent{Type: "message"},
		ID:       id,
		Channel:  channel,
		User:     "",
		Text:     text,
	}
	return websocket.JSON.Send(s.websocket, m)
}

// Ping sends a ping request
func (s *Connection) Ping() error {
	id := nextID()
	m := RtmPingPong{
		RtmEvent:  RtmEvent{Type: "ping"},
		ID:        id,
		Timestamp: time.Now(),
	}
	return websocket.JSON.Send(s.websocket, m)
}
