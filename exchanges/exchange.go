// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package exchanges

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SideType string

const (
	SideTypeSell SideType = "sell"
	SideTypeBuy  SideType = "buy"
)

// TickerEvent struct
type TickerEvent struct {
	Product Product   `json:"product"`
	Price   float64   `json:"price"`
	Side    SideType  `json:"side"`
	Time    time.Time `json:"time"`
	Size    float64   `json:"size"`
}

// OrderEvent struct
type OrderEvent struct {
	Side  SideType
	Price float64
}

// Product struct
type Product struct {
	From string
	To   string
}

// NewProductFromString create product object
func NewProductFromString(product string) Product {
	part := strings.Split(product, "-")

	return Product{
		From: part[0],
		To:   part[1],
	}
}

// NewProduct create product object
func NewProduct(from string, to string) Product {
	return Product{
		From: from,
		To:   to,
	}
}

func (p Product) String() string {
	return fmt.Sprintf(`%s-%s`, p.From, p.To)
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (p Product) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It support string
// and null input.
func (p *Product) UnmarshalJSON(data []byte) error {
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

// ExchangeProvider interface
type ExchangeProvider interface {
	Name() string

	Ticker() TickerProvider
	// Order(from string, to string) (<-chan *OrderEvent, error)
}

// TickerProvider interface
type TickerProvider interface {
	Subscribe(products ...Product) error
	Unsubscribe(products ...Product) error
	Channel() <-chan *TickerEvent
}
