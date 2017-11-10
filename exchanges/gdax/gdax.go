// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gdax

import (
	"strconv"

	"github.com/euskadi31/cryptotrader/exchanges"
	gdaxclient "github.com/preichenberger/go-gdax"
	"github.com/rs/zerolog/log"
)

// GDAX struct
type GDAX struct {
	client *gdaxclient.Client
	ws     *WebSocketClient
}

// NewGDAX Exchange
func NewGDAX() (*GDAX, error) {
	e := &GDAX{
		client: gdaxclient.NewClient("", "", ""),
		ws:     NewWebSocketClient(),
	}

	if err := e.ws.Connect(); err != nil {
		return nil, err
	}

	return e, nil
}

// Ticker channel
func (e *GDAX) Ticker(from string, to string) (<-chan *exchanges.TickerEvent, error) {
	log.Debug().Msg("Ticker called")

	out := make(chan *exchanges.TickerEvent)

	if err := e.ws.Subscribe(&WebSocketChannel{
		Name: WebSocketChannelTypeTicker,
		Products: []*WebSocketProduct{
			&WebSocketProduct{
				From: from,
				To:   to,
			},
		},
	}); err != nil {
		return out, err
	}

	go func() {
		for {
			select {
			case t := <-e.ws.Ticker:
				price, err := strconv.ParseFloat(t.Price, 64)
				if err != nil {
					log.Error().Err(err).Msg("")
				}

				size, err := strconv.ParseFloat(t.LastSize, 64)
				if err != nil {
					log.Error().Err(err).Msg("")
				}

				side := exchanges.SideTypeBuy

				if t.Side == "sell" {
					side = exchanges.SideTypeSell
				}

				out <- &exchanges.TickerEvent{
					Price: price,
					Time:  t.Time.Time(),
					Side:  side,
					Size:  size,
				}
			}
		}
	}()

	return out, nil
}
