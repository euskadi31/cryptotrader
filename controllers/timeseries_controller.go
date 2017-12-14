// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/euskadi31/cryptotrader/trader"

	"github.com/euskadi31/go-server"
)

// TimeseriesController struct
type TimeseriesController struct {
	engine *trader.Engine
}

// NewTimeseriesController constructor
func NewTimeseriesController(engine *trader.Engine) *TimeseriesController {
	return &TimeseriesController{
		engine: engine,
	}
}

// Mount implements server.Controller
func (c *TimeseriesController) Mount(r *server.Router) {
	r.AddRouteFunc("/api/v1/timeseries/{provider:[a-z]+}/{from:[a-z]+}-{to:[a-z]+}", c.GetTimeseriesHandler).Methods(http.MethodGet)
}

// GetTimeseriesHandler endpoint
func (c *TimeseriesController) GetTimeseriesHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	key := fmt.Sprintf("%s-%s-%s", params["provider"], strings.ToUpper(params["from"]), strings.ToUpper(params["to"]))

	ts, err := c.engine.GetTimeserie(key)
	if err != nil {
		server.FailureFromError(w, http.StatusNotFound, err)

		return
	}

	server.JSON(w, http.StatusOK, ts.All())
}
