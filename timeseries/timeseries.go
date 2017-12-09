// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package timeseries

import (
	"math"
	"sync"

	"github.com/toolsparty/regression"
)

// TrendType type
type TrendType int

// TrendType enum
const (
	TrendTypeDecreasing TrendType = -1
	TrendTypeNeutral    TrendType = 0
	TrendTypeIncreasing TrendType = 1
)

// DataPoint struct
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
func (ts Timeseries) Keys() []int64 {
	return ts.times
}

// Values slice
func (ts Timeseries) Values() []float64 {
	return ts.values
}

// MaxValue of Timeseries
func (ts Timeseries) MaxValue() float64 {
	max := float64(0)

	ts.RLock()
	for _, v := range ts.values {
		max = math.Max(max, v)
	}
	ts.RUnlock()

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

// GetLatestValues by size
func (ts Timeseries) GetLatestValues(size int) []float64 {
	length := len(ts.values)

	if length < size {
		return ts.values
	}

	ts.RLock()
	values := ts.values[length-size:]
	ts.RUnlock()

	return values
}

// GetTrending for timeseries
func (ts Timeseries) GetTrending(size int) (TrendType, error) {
	reg, err := regression.NewLinear([]float64{}, ts.GetLatestValues(size))
	if err != nil {
		return TrendTypeNeutral, err
	}

	theta := reg.GetTheta()

	if theta > 0 {
		return TrendTypeIncreasing, nil
	} else if theta < 0 {
		return TrendTypeDecreasing, nil
	}

	return TrendTypeNeutral, nil
}
