// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"testing"

	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {

	marketPrice := float64(14199)

	o := &Order{
		Provider:  "gdax",
		Side:      exchanges.SideTypeBuy,
		ProductID: "BTC-EUR",
		Size:      0.04468611,
		Price:     150,
	}

	assert.Equal(t, 3356.7477679305716, o.GetBuyingMarketPrice())
	assert.Equal(t, 634.49807589, o.GetCurrentPrice(marketPrice))
	assert.Equal(t, 484.49807589, o.GetMarginInCurrency(marketPrice))
	assert.Equal(t, 23.640733628639847, o.GetMarginInPercent(marketPrice))
}

func BenchmarkOrderGetBuyingMarketPrice(b *testing.B) {
	o := &Order{
		Provider:  "gdax",
		Side:      exchanges.SideTypeBuy,
		ProductID: "BTC-EUR",
		Size:      0.04468611,
		Price:     150,
	}

	for n := 0; n < b.N; n++ {
		o.GetBuyingMarketPrice()
	}
}

func BenchmarkOrderGetCurrentPrice(b *testing.B) {
	marketPrice := float64(14199)

	o := &Order{
		Provider:  "gdax",
		Side:      exchanges.SideTypeBuy,
		ProductID: "BTC-EUR",
		Size:      0.04468611,
		Price:     150,
	}

	for n := 0; n < b.N; n++ {
		o.GetCurrentPrice(marketPrice)
	}
}

func BenchmarkOrderGetMarginInCurrency(b *testing.B) {
	marketPrice := float64(14199)

	o := &Order{
		Provider:  "gdax",
		Side:      exchanges.SideTypeBuy,
		ProductID: "BTC-EUR",
		Size:      0.04468611,
		Price:     150,
	}

	for n := 0; n < b.N; n++ {
		o.GetMarginInCurrency(marketPrice)
	}
}

func BenchmarkOrderGetMarginInPercent(b *testing.B) {
	marketPrice := float64(14199)

	o := &Order{
		Provider:  "gdax",
		Side:      exchanges.SideTypeBuy,
		ProductID: "BTC-EUR",
		Size:      0.04468611,
		Price:     150,
	}

	for n := 0; n < b.N; n++ {
		o.GetMarginInPercent(marketPrice)
	}
}
