// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"github.com/euskadi31/go-std"
)

// CampaignState type
type CampaignState string

// CampaignState enum
const (
	CampaignStateSelling CampaignState = "selling"
	CampaignStateBuying  CampaignState = "buying"
)

// Campaign struct
type Campaign struct {
	ID          int           `db:"id,increment"`
	Provider    string        `db:"provider" json:"provider"`
	ProviderRef string        `db:"provider_ref" json:"provider_ref"`
	ProductID   string        `db:"product_id" json:"product_id"`
	Volume      float64       `db:"volume" json:"volume"`
	BuyLimit    float64       `db:"buy_limit" json:"buy_limit"`
	SellLimit   float64       `db:"sell_limit" json:"sell_limit"`
	CreatedAt   std.DateTime  `db:"created_at" json:"created_at"`
	UpdatedAt   std.DateTime  `db:"updated_at" json:"updated_at"`
	DeletedAt   std.DateTime  `db:"deleted_at" json:"deleted_at"`
	Orders      []*Order      `db:"-" json:"orders"`
	State       CampaignState `db:"state" json:"state"`
}

// AddOrder to Campaign
func (c *Campaign) AddOrder(order *Order) {
	c.Orders = append(c.Orders, order)
}

// IsSelling returns true if campaign is state selling
func (c *Campaign) IsSelling() bool {
	return c.State == CampaignStateSelling
}

// IsBuying returns true if campaign is state buying
func (c *Campaign) IsBuying() bool {
	return c.State == CampaignStateBuying
}
