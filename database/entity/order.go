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
	Provider  string             `db:"provider" json:"provider"`
	ID        string             `db:"id" json:"id"`
	TradeID   string             `db:"trade_id" json:"trade_id"`
	Side      exchanges.SideType `db:"side" json:"side"`
	Size      float64            `db:"size" json:"size"`
	ProductID string             `db:"product_id" json:"product_id"`
	Price     float64            `db:"price" json:"price"`
	CreatedAt std.DateTime       `db:"created_at" json:"created_at"`
	UpdatedAt std.DateTime       `db:"updated_at" json:"updated_at"`
	DeletedAt std.DateTime       `db:"deleted_at" json:"deleted_at"`
}
