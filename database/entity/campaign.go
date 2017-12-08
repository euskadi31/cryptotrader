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
	CampaignStateSell    CampaignState = "sell"
	CampaignStateSelling CampaignState = "selling"
	CampaignStateBuy     CampaignState = "buy"
	CampaignStateBuying  CampaignState = "buying"
)

// Campaign struct
type Campaign struct {
	ID            int           `storm:"id,increment"`
	Provider      string        `storm:"index" json:"provider"`
	ProviderRef   string        `json:"provider_ref"`
	ProductID     string        `storm:"index" json:"product_id"`
	Volume        float64       `json:"volume"`
	BuyLimit      float64       `json:"buy_limit"`
	SellLimit     float64       `json:"sell_limit"`
	SellLimitUnit string        `json:"sell_limit_unit"`
	CreatedAt     std.DateTime  `json:"created_at"`
	UpdatedAt     std.DateTime  `json:"updated_at"`
	BuyOrder      *Order        `json:"buy_order"`
	SellOrder     *Order        `json:"sell_order"`
	Orders        []*Order      `json:"orders"`
	State         CampaignState `storm:"index" json:"state"`
}

// AddOrder to Campaign
func (c *Campaign) AddOrder(order *Order) {
	c.Orders = append(c.Orders, order)
}

// IsState check if this state is eq to state param
func (c Campaign) IsState(state CampaignState) bool {
	return c.State == state
}

// IsSelling returns true if campaign is state selling
func (c *Campaign) IsSelling() bool {
	return c.State == CampaignStateSelling
}

// IsBuying returns true if campaign is state buying
func (c *Campaign) IsBuying() bool {
	return c.State == CampaignStateBuying
}
