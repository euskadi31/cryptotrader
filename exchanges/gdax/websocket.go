// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gdax

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WebSocketEventType string

const (
	WebSocketEventTypeSubscribe           WebSocketEventType = "subscribe"
	WebSocketEventTypeSubscriptions       WebSocketEventType = "subscriptions"
	WebSocketEventTypeError               WebSocketEventType = "error"
	WebSocketEventTypeUnsubscribe         WebSocketEventType = "unsubscribe"
	WebSocketEventTypeOpen                WebSocketEventType = "open"
	WebSocketEventTypeReceived            WebSocketEventType = "received"
	WebSocketEventTypeDone                WebSocketEventType = "done"
	WebSocketEventTypeMatch               WebSocketEventType = "match"
	WebSocketEventTypeChange              WebSocketEventType = "change"
	WebSocketEventTypeActivate            WebSocketEventType = "activate"
	WebSocketEventTypeHeartbeat           WebSocketEventType = "heartbeat"
	WebSocketEventTypeTicker              WebSocketEventType = "ticker"
	WebSocketEventTypeSnapshot            WebSocketEventType = "snapshot"
	WebSocketEventTypeLevel2Update        WebSocketEventType = "l2update"
	WebSocketEventTypeMarginProfileUpdate WebSocketEventType = "margin_profile_update"
)

// WebSocketEvent struct
type WebSocketEvent struct {
	Type WebSocketEventType `json:"type"`
}

// WebSocketProduct struct
type WebSocketProduct struct {
	From string
	To   string
}

// NewWebSocketProduct create product object
func NewWebSocketProduct(from string, to string) *WebSocketProduct {
	return &WebSocketProduct{
		From: from,
		To:   to,
	}
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (p WebSocketProduct) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%s-%s", p.From, p.To)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It support string
// and null input.
func (p *WebSocketProduct) UnmarshalJSON(data []byte) error {
	var v string

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	part := strings.Split(v, "-")

	p.From = part[0]
	p.To = part[1]

	return nil
}

// WebSocketChannel struct
type WebSocketChannel struct {
	Name     string              `json:"name"`
	Products []*WebSocketProduct `json:"product_ids,omitempty"`
}

// WebSocketSubscribeRequest struct
type WebSocketSubscribeRequest struct {
	*WebSocketEvent
	Channels []*WebSocketChannel `json:"channels"`
}

// WebSocketTickerResponse struct
type WebSocketTickerResponse struct {
	*WebSocketEvent
	TradeID  int               `json:"trade_id"`
	Sequence int               `json:"sequence"`
	Product  *WebSocketProduct `json:"product_id"`
	Price    string            `json:"price"`
	Side     string            `json:"side"`
	LastSize string            `json:"last_size"`
	BestBid  string            `json:"best_bid"`
	BestAsk  string            `json:"best_ask"`
}

// WebSocketClient struct
type WebSocketClient struct {
	api    string
	ws     *websocket.Conn
	Ticker chan *WebSocketTickerResponse
}

// NewWebSocketClient constructor
func NewWebSocketClient() *WebSocketClient {
	ws := &WebSocketClient{
		api:    "wss://ws-feed.gdax.com",
		Ticker: make(chan *WebSocketTickerResponse, 100),
	}

	go ws.receiver()

	return ws
}

// Connect to websocket server
func (c *WebSocketClient) Connect() error {
	u, err := url.Parse(c.api)
	if err != nil {
		return err
	}

	log.Info().Msgf("GDAX WebSocket: connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.ws = ws

	return nil
}

// Subscribe to channel
func (c *WebSocketClient) Subscribe(channels ...*WebSocketChannel) error {

	e := &WebSocketSubscribeRequest{
		WebSocketEvent: &WebSocketEvent{
			Type: WebSocketEventTypeSubscribe,
		},
		Channels: channels,
	}

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return c.ws.WriteMessage(websocket.TextMessage, b)
}

// Unsubscribe to channel
func (c *WebSocketClient) Unsubscribe(channels ...*WebSocketChannel) error {
	e := &WebSocketSubscribeRequest{
		WebSocketEvent: &WebSocketEvent{
			Type: WebSocketEventTypeUnsubscribe,
		},
		Channels: channels,
	}

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return c.ws.WriteMessage(websocket.TextMessage, b)
}

func (c *WebSocketClient) receiver() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("")

			continue
		}

		var evtType WebSocketEvent

		if err := json.Unmarshal(message, &evtType); err != nil {
			log.Error().Err(err).Msg("")

			continue
		}

		if err := c.processEvent(evtType.Type, message); err != nil {
			log.Error().Err(err).Msg("")

			continue
		}
	}
}

func (c *WebSocketClient) processEvent(event WebSocketEventType, data []byte) error {
	switch event {
	case WebSocketEventTypeError:
		// @TODO parse error message
		return nil
	case WebSocketEventTypeTicker:
		v := &WebSocketTickerResponse{}

		if err := json.Unmarshal(data, v); err != nil {
			return err
		}

		c.Ticker <- v

		return nil
	default:
		return errors.New("Bad event type")
	}
}
