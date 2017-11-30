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
	Provider  string             `storm:"provider" json:"provider"`
	ID        string             `storm:"id" json:"id"`
	TradeID   string             `storm:"trade_id" json:"trade_id"`
	Side      exchanges.SideType `storm:"side" json:"side"`
	Size      float64            `storm:"size" json:"size"`
	ProductID string             `storm:"product_id" json:"product_id"`
	Price     float64            `storm:"price" json:"price"`
	CreatedAt std.DateTime       `storm:"created_at" json:"created_at"`
	UpdatedAt std.DateTime       `storm:"updated_at" json:"updated_at"`
	DeletedAt std.DateTime       `storm:"deleted_at" json:"deleted_at"`
}
