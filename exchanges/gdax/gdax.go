// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gdax

import (
	"strconv"

	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/rs/zerolog/log"
)

// GDAX struct
type GDAX struct {
	ws *WebSocketClient
}

// NewGDAX Exchange
func NewGDAX() (*GDAX, error) {
	e := &GDAX{
		ws: NewWebSocketClient(),
	}

	if err := e.ws.Connect(); err != nil {
		return nil, err
	}

	return e, nil
}

// Ticker channel
func (e *GDAX) Ticker() (<-chan *exchanges.TickerEvent, error) {
	out := make(chan *exchanges.TickerEvent)

	e.ws.Subscribe(&WebSocketChannel{
		Name: "ticker",
		Products: []*WebSocketProduct{
			&WebSocketProduct{
				From: "BTC",
				To:   "EUR",
			},
		},
	})

	go func() {
		for {
			select {
			case t := <-e.ws.Ticker:
				f, err := strconv.ParseFloat(t.Price, 64)
				if err != nil {
					log.Error().Err(err).Msg("")
				}

				out <- &exchanges.TickerEvent{
					Price: f,
				}
			}
		}
	}()

	return out, nil
}
