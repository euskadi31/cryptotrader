// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package timeseries

import (
	"math"
	"sync"
)

type DataPoint struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

// Timeseries struct
type Timeseries struct {
	*sync.RWMutex
	times  []int64
	values []float64
}

// New Timeseries
func New() *Timeseries {
	return &Timeseries{
		RWMutex: &sync.RWMutex{},
		times:   []int64{},
		values:  []float64{},
	}
}

// Add DataPoint to Timeseries
func (ts *Timeseries) Add(t int64, v float64) {
	ts.Lock()
	ts.times = append(ts.times, t)
	ts.values = append(ts.values, v)
	ts.Unlock()
}

// Keys slice
func (ts *Timeseries) Keys() []int64 {
	return ts.times
}

// Values slice
func (ts *Timeseries) Values() []float64 {
	return ts.values
}

// MaxValue of Timeseries
func (ts *Timeseries) MaxValue() float64 {
	max := float64(0)

	for _, v := range ts.values {
		max = math.Max(max, v)
	}

	return max
}

// All DataPoint in Timeseries
func (ts *Timeseries) All() []*DataPoint {
	datas := []*DataPoint{}
	ts.RLock()
	for i, v := range ts.values {
		t := ts.times[i]

		datas = append(datas, &DataPoint{
			Time:  t,
			Value: v,
		})
	}
	ts.RUnlock()

	return datas
}
