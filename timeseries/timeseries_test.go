// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package timeseries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeseriesWithSize(t *testing.T) {
	ts := New(6)

	ts.Add(1, 15)
	ts.Add(2, 23)
	ts.Add(3, 45)
	ts.Add(4, 21)
	ts.Add(5, 11)
	ts.Add(6, 9)
	ts.Add(7, 10)

	assert.Equal(t, 6, ts.Size())
	assert.Equal(t, []int64{2, 3, 4, 5, 6, 7}, ts.Keys())
	assert.Equal(t, []float64{23, 45, 21, 11, 9, 10}, ts.Values())
	assert.Equal(t, float64(45), ts.MaxValue())
	assert.Equal(t, []float64{11, 9, 10}, ts.GetLatestValues(3))
	assert.Equal(t, []float64{23, 45, 21, 11, 9, 10}, ts.GetLatestValues(10))

	assert.Equal(t, []*DataPoint{
		&DataPoint{
			Time:  2,
			Value: 23,
		},
		&DataPoint{
			Time:  3,
			Value: 45,
		},
		&DataPoint{
			Time:  4,
			Value: 21,
		},
		&DataPoint{
			Time:  5,
			Value: 11,
		},
		&DataPoint{
			Time:  6,
			Value: 9,
		},
		&DataPoint{
			Time:  7,
			Value: 10,
		},
	}, ts.All())
}

func TestTimeseries(t *testing.T) {
	ts := New(6)

	ts.Add(1, 15)
	ts.Add(2, 23)
	ts.Add(3, 45)
	ts.Add(4, 21)
	ts.Add(5, 11)
	ts.Add(6, 9)

	assert.Equal(t, []int64{1, 2, 3, 4, 5, 6}, ts.Keys())
	assert.Equal(t, []float64{15, 23, 45, 21, 11, 9}, ts.Values())
	assert.Equal(t, float64(45), ts.MaxValue())
	assert.Equal(t, []float64{21, 11, 9}, ts.GetLatestValues(3))
	assert.Equal(t, []float64{15, 23, 45, 21, 11, 9}, ts.GetLatestValues(10))

	assert.Equal(t, []*DataPoint{
		&DataPoint{
			Time:  1,
			Value: 15,
		},
		&DataPoint{
			Time:  2,
			Value: 23,
		},
		&DataPoint{
			Time:  3,
			Value: 45,
		},
		&DataPoint{
			Time:  4,
			Value: 21,
		},
		&DataPoint{
			Time:  5,
			Value: 11,
		},
		&DataPoint{
			Time:  6,
			Value: 9,
		},
	}, ts.All())
}

func TestTimeseriesTrendingIncreasing(t *testing.T) {
	ts := New(6)

	ts.Add(1, 15)
	ts.Add(2, 23)
	ts.Add(3, 45)
	ts.Add(4, 43)
	ts.Add(5, 46)
	ts.Add(6, 10)

	trend, err := ts.GetTrending(10)
	assert.NoError(t, err)

	assert.Equal(t, TrendTypeIncreasing, trend)
}

func TestTimeseriesTrendingDecreasing(t *testing.T) {
	ts := New(6)

	ts.Add(1, 38)
	ts.Add(2, 39)
	ts.Add(3, 40)
	ts.Add(4, 25)
	ts.Add(5, 20)
	ts.Add(6, 10)

	trend, err := ts.GetTrending(10)
	assert.NoError(t, err)

	assert.Equal(t, TrendTypeDecreasing, trend)
}

func TestTimeseriesTrendingNeutral(t *testing.T) {
	ts := New(6)

	ts.Add(1, 10)
	ts.Add(2, 10)
	ts.Add(3, 10)
	ts.Add(4, 10)
	ts.Add(5, 10)
	ts.Add(6, 10)

	trend, err := ts.GetTrending(10)
	assert.NoError(t, err)

	assert.Equal(t, TrendTypeNeutral, trend)
}

func BenchmarkTimeseries(b *testing.B) {
	b.ReportAllocs()
	ts := New(b.N + 20)

	for n := 0; n < (b.N * 2); n++ {
		ts.Add(int64(n), float64(n*10))
	}

	for n := 0; n < b.N; n++ {
		ts.GetTrending(n + 10)
	}
}

func BenchmarkTimeseriesAddWithSize(b *testing.B) {
	b.ReportAllocs()
	ts := New(4)

	for n := 0; n < b.N; n++ {
		ts.Add(int64(n), float64(n))
	}
}

func BenchmarkTimeseriesAdd(b *testing.B) {
	b.ReportAllocs()
	ts := New(b.N + 20)

	for n := 0; n < b.N; n++ {
		ts.Add(int64(n), float64(n*10))
	}
}
