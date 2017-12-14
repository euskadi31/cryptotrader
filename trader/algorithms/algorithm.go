// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algorithms

import (
	"encoding/json"

	"github.com/euskadi31/cryptotrader/database/entity"
	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/euskadi31/cryptotrader/timeseries"
)

// Algorithm interface
type Algorithm interface {
	json.Marshaler

	Name() string

	Options() Options

	Buy(event *exchanges.TickerEvent, campaign *entity.Campaign, ts *timeseries.Timeseries)

	Sell(event *exchanges.TickerEvent, campaign *entity.Campaign, ts *timeseries.Timeseries)
}
