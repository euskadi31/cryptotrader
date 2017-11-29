// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package services

import "github.com/euskadi31/cryptotrader/exchanges"

// Ticker event
type Ticker struct {
	*exchanges.TickerEvent
	Product string
}

// TickerObserver alias
type TickerObserver func(event *Ticker)

// TickerService struc
type TickerService struct {
	eventsCh  chan *Ticker
	observers map[string][]TickerObserver
}

// NewTickerService func
func NewTickerService() *TickerService {
	s := &TickerService{
		eventsCh:  make(chan *Ticker, 100),
		observers: make(map[string][]TickerObserver),
	}

	go s.looper()

	return s
}

// AddObserver to service
func (s *TickerService) AddObserver(channel string, observer TickerObserver) {
	s.observers[channel] = append(s.observers[channel], observer)
}

func (s *TickerService) looper() {
	for {
		select {
		case e := <-s.eventsCh:
			observers, ok := s.observers[e.Product]
			if ok == false {
				continue
			}

			for _, observer := range observers {
				observer(e)
			}
		}
	}
}

// Send Ticker event to service
func (s *TickerService) Send(event *Ticker) {
	s.eventsCh <- event
}
