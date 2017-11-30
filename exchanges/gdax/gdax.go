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

// Name of provider
func (e GDAX) Name() string {
	return "gdax"
}

// Ticker channel
func (e *GDAX) Ticker() exchanges.TickerProvider {
	return &Ticker{
		ws: e.ws,
	}
}

// Ticker struct
type Ticker struct {
	ws *WebSocketClient
}

func (t Ticker) convertProduct(products []exchanges.Product) []*WebSocketProduct {
	sp := []*WebSocketProduct{}

	for _, p := range products {
		sp = append(sp, &WebSocketProduct{
			From: p.From,
			To:   p.To,
		})
	}

	return sp
}

// Subscribe to product
func (t *Ticker) Subscribe(products ...exchanges.Product) error {
	sp := []*WebSocketProduct{}

	for _, p := range products {
		sp = append(sp, &WebSocketProduct{
			From: p.From,
			To:   p.To,
		})
	}

	if err := t.ws.Subscribe(&WebSocketChannel{
		Name:     WebSocketChannelTypeTicker,
		Products: t.convertProduct(products),
	}); err != nil {
		return err
	}

	return nil
}

// Unsubscribe to product
func (t *Ticker) Unsubscribe(products ...exchanges.Product) error {
	sp := []*WebSocketProduct{}

	for _, p := range products {
		sp = append(sp, &WebSocketProduct{
			From: p.From,
			To:   p.To,
		})
	}

	if err := t.ws.Unsubscribe(&WebSocketChannel{
		Name:     WebSocketChannelTypeTicker,
		Products: t.convertProduct(products),
	}); err != nil {
		return err
	}

	return nil
}

// Channel TickerEvent
func (t *Ticker) Channel() <-chan *exchanges.TickerEvent {
	out := make(chan *exchanges.TickerEvent)

	go func() {
		for {
			select {
			case t := <-t.ws.Ticker:
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
					Product: exchanges.NewProduct(t.Product.From, t.Product.To),
					Price:   price,
					Time:    t.Time.Time(),
					Side:    side,
					Size:    size,
				}
			}
		}
	}()

	return out
}
