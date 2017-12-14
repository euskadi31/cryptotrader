// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algorithms

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyCampaign struct {
	Options Options `json:"options"`
}

var jsonData = []byte(`{
	"options": {
		"trend.minimal_history_size": 150,
		"trend.max_price": 150.50,
		"trend.type": "percent",
		"trend.activate": true
	}
}`)

func TestOptionsNotInit(t *testing.T) {
	var o Options

	assert.Equal(t, 0, o.GetInt("Trend.minimal_history_size"))
	assert.Equal(t, 0.0, o.GetFloat("Trend.max_price"))
	assert.Equal(t, false, o.GetBool("Trend.activate"))
	assert.Equal(t, "", o.GetString("Trend.type"))
}

func TestOptionsEmpty(t *testing.T) {
	o := &Options{}

	assert.Equal(t, 0, o.GetInt("Trend.minimal_history_size"))
	assert.Equal(t, 0.0, o.GetFloat("Trend.max_price"))
	assert.Equal(t, false, o.GetBool("Trend.activate"))
	assert.Equal(t, "", o.GetString("Trend.type"))
}

func TestOptionsDefault(t *testing.T) {
	campaign := &MyCampaign{
		Options: Options{
			"trend.max_price": 140.0,
			"trend.long_size": 10,
			"status":          false,
		},
	}

	err := json.Unmarshal(jsonData, campaign)
	assert.NoError(t, err)

	assert.Equal(t, 10, campaign.Options.GetInt("trend.long_size"))
	assert.Equal(t, 150, campaign.Options.GetInt("Trend.minimal_history_size"))
	assert.Equal(t, 150.50, campaign.Options.GetFloat("Trend.max_price"))
	assert.Equal(t, true, campaign.Options.GetBool("Trend.activate"))
	assert.Equal(t, "percent", campaign.Options.GetString("Trend.type"))

	assert.Equal(t, false, campaign.Options.GetBool("status"))
}

func TestOptionsFromJSON(t *testing.T) {
	campaign := &MyCampaign{}

	err := json.Unmarshal(jsonData, campaign)
	assert.NoError(t, err)

	campaign.Options.Set("foo", "bar")

	assert.Equal(t, 150, campaign.Options.GetInt("Trend.minimal_history_size"))
	assert.Equal(t, 150.50, campaign.Options.GetFloat("Trend.max_price"))
	assert.Equal(t, true, campaign.Options.GetBool("Trend.activate"))
	assert.Equal(t, "percent", campaign.Options.GetString("Trend.type"))

	assert.Equal(t, "bar", campaign.Options.GetString("foo"))
}

func TestOptionsMerge(t *testing.T) {
	defaultOptions := Options{
		"trend.max_price": 140.0,
		"status":          false,
	}

	options := Options{
		"trend.max_price": 150.0,
	}

	defaultOptions.Merge(options)

	assert.Equal(t, 150.00, defaultOptions.GetFloat("Trend.max_price"))
	assert.Equal(t, false, defaultOptions.GetBool("status"))
}
