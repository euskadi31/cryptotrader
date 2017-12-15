// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algorithms

import (
	"encoding/json"
	"testing"

	"github.com/euskadi31/cryptotrader/database/entity"
	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/euskadi31/cryptotrader/timeseries"
	"github.com/stretchr/testify/assert"
)

type MyAlgo struct {
}

// Name implements Algorithm interface
func (a MyAlgo) Name() string {
	return "my-algo"
}

// Options implements Algorithm interface
func (a MyAlgo) Options() Options {
	return Options{}
}

// MarshalJSON implements json.Marshaler.
func (a MyAlgo) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Options())
}

// Buy implements Algorithm interface
func (a *MyAlgo) Buy(event *exchanges.TickerEvent, campaign *entity.Campaign, ts *timeseries.Timeseries) {

}

// Sell implements Algorithm interface
func (a *MyAlgo) Sell(event *exchanges.TickerEvent, campaign *entity.Campaign, ts *timeseries.Timeseries) {

}

func TestManagerNotInit(t *testing.T) {
	var m Manager

	algo, err := m.Get("test-not-found")
	assert.Error(t, err)
	assert.EqualError(t, err, ErrManagerNotInit.Error())
	assert.Nil(t, algo)
}

func TestManager(t *testing.T) {
	m := NewManager()

	assert.False(t, m.Has("test-not-found"))

	algo, err := m.Get("test-not-found")
	assert.Error(t, err)
	assert.EqualError(t, err, ErrAlgorithmNotFound.Error())
	assert.Nil(t, algo)

	a := &MyAlgo{}
	m.Add(a)

	algo, err = m.Get("my-algo")
	assert.NoError(t, err)
	assert.Equal(t, a, algo)
}
