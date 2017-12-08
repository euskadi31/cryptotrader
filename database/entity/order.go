// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/euskadi31/go-std"
)

// Order struct
type Order struct {
	Provider  string             `json:"provider"`
	ID        int                `storm:"id,increment" json:"id"`
	TradeID   string             `json:"trade_id"`
	Side      exchanges.SideType `json:"side"`
	Size      float64            `json:"size"`
	ProductID string             `json:"product_id"`
	Price     float64            `json:"price"`
	CreatedAt std.DateTime       `json:"created_at"`
	UpdatedAt std.DateTime       `json:"updated_at"`
	DeletedAt std.DateTime       `json:"deleted_at"`
}

// GetBuyingMarketPrice func
func (o Order) GetBuyingMarketPrice() float64 {
	return o.Price / o.Size
}

// GetCurrentPrice from market price
func (o Order) GetCurrentPrice(marketPrice float64) float64 {
	return marketPrice * o.Size
}

// GetMarginInCurrency from market price
func (o Order) GetMarginInCurrency(marketPrice float64) float64 {
	return o.GetCurrentPrice(marketPrice) - o.Price
}

// GetMarginInPercent from market price
func (o Order) GetMarginInPercent(marketPrice float64) float64 {
	return 100 * (o.Price / o.GetCurrentPrice(marketPrice))
}
