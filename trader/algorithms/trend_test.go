// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algorithms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrendName(t *testing.T) {
	algo := NewTrend(nil, nil)

	assert.Equal(t, "trend", algo.Name())
}

func TestTrendOptions(t *testing.T) {
	algo := NewTrend(nil, nil)

	assert.Equal(t, Options{
		"trend.selling.long_trend_size":  150,
		"trend.selling.short_trend_size": 10,
	}, algo.Options())
}
