// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gdax

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	wsBase = "wss://ws-feed.gdax.com"
)

type wsEvent struct {
	Type string `json:"type"`
}

// GDAX struct
type GDAX struct {
}

func (e *GDAX) processEvent(event string, data []byte) (*exchanges.TickerEvent, error) {

	switch event {
	case "error":
	case "ticker":
		return &exchanges.TickerEvent{}, nil
	default:
		return nil, errors.New("Bad event type")
	}

	return nil, nil
}

func (e *GDAX) Ticker() (<-chan *exchanges.TickerEvent, error) {
	out := make(chan *exchanges.TickerEvent)

	u, err := url.Parse(wsBase)
	if err != nil {
		return out, err
	}

	log.Info().Msgf("GDAX: connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return out, err
	}

	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Error().Err(err).Msg("")

				continue
			}

			var evtType wsEvent

			if err := json.Unmarshal(message, &evtType); err != nil {
				log.Error().Err(err).Msg("")

				continue
			}

			evt, err := e.processEvent(evtType.Type, message)
			if err != nil {
				log.Error().Err(err).Msg("")

				continue
			}

			out <- evt
		}
	}()

	// @TODO: subscribe to ticker channel

	return out, nil
}
