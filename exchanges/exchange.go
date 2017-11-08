// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package exchanges

import (
	"time"
)

type SideType string

const (
	SideTypeSell SideType = "sell"
	SideTypeBuy  SideType = "buy"
)

// TickerEvent struct
type TickerEvent struct {
	Price float64
	Side  SideType
	Time  time.Time
	Size  float64
}

// OrderEvent struct
type OrderEvent struct {
	Side  SideType
	Price float64
}

// ExchangeProvider interface
type ExchangeProvider interface {
	Ticker(from string, to string) (<-chan *TickerEvent, error)
	// Order(from string, to string) (<-chan *OrderEvent, error)
}
