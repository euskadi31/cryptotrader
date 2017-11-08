// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/euskadi31/cryptotrader/timeseries"
	"github.com/euskadi31/go-server"
)

// TimeseriesController struct
type TimeseriesController struct {
	ts *timeseries.Timeseries
}

// NewTimeseriesController constructor
func NewTimeseriesController(ts *timeseries.Timeseries) *TimeseriesController {
	return &TimeseriesController{
		ts: ts,
	}
}

// Mount implements server.Controller
func (c *TimeseriesController) Mount(r *server.Router) {
	r.AddRouteFunc("/api/v1/timeseries", c.GetTimeseriesHandler).Methods(http.MethodGet)
}

// GetTimeseriesHandler endpoint
func (c *TimeseriesController) GetTimeseriesHandler(w http.ResponseWriter, r *http.Request) {
	server.JSON(w, http.StatusOK, c.ts.All())
}
