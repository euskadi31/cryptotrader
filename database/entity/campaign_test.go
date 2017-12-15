// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCampaign(t *testing.T) {
	c := &Campaign{
		State: CampaignStateBuy,
	}

	assert.True(t, c.IsState(CampaignStateBuy))
	assert.False(t, c.IsSelling())
	assert.False(t, c.IsBuying())

	c.State = CampaignStateSelling

	assert.True(t, c.IsState(CampaignStateSelling))
	assert.True(t, c.IsSelling())
	assert.False(t, c.IsBuying())

	c.State = CampaignStateBuying

	assert.True(t, c.IsState(CampaignStateBuying))
	assert.False(t, c.IsSelling())
	assert.True(t, c.IsBuying())

	assert.Equal(t, 0, len(c.Orders))

	c.AddOrder(&Order{})

	assert.Equal(t, 1, len(c.Orders))
}
