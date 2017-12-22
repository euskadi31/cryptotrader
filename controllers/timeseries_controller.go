// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/euskadi31/cryptotrader/trader"
	"github.com/euskadi31/go-eventemitter"
	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-sse"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// TimeseriesController struct
type TimeseriesController struct {
	engine  *trader.Engine
	emitter eventemitter.EventEmitter
}

// NewTimeseriesController constructor
func NewTimeseriesController(engine *trader.Engine, emitter eventemitter.EventEmitter) *TimeseriesController {
	return &TimeseriesController{
		engine:  engine,
		emitter: emitter,
	}
}

// Mount implements server.Controller
func (c *TimeseriesController) Mount(r *server.Router) {
	events := sse.NewServer(c.GetTimeseriesEventHandler)
	events.SetRetry(time.Second * 5)

	r.AddRouteFunc("/api/v1/timeseries/{provider:[a-z]+}/{from:[a-z]+}-{to:[a-z]+}", c.GetTimeseriesHandler).Methods(http.MethodGet)
	r.AddRoute("/api/v1/timeseries/{provider:[a-z]+}/{from:[a-z]+}-{to:[a-z]+}/events", events).Methods(http.MethodGet)
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

// GetTimeseriesEventHandler endpoint
func (c *TimeseriesController) GetTimeseriesEventHandler(rw sse.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	key := fmt.Sprintf("%s-%s-%s", params["provider"], strings.ToUpper(params["from"]), strings.ToUpper(params["to"]))
	eventKey := fmt.Sprintf("ticker-%s", key)

	// recovery from timeseries
	// @TODO: implements this
	lastID := r.Header.Get("Last-Event-ID")
	if lastID != "" {
		log.Debug().Msgf("Recovery with ID: %s", lastID)
	}

	listener := func(event *exchanges.TickerEvent) {
		b, err := json.Marshal(event)
		if err != nil {
			log.Error().Err(err).Msg("Marshal TickerEvent failed")

			return
		}

		rw.Send(&sse.MessageEvent{
			ID:   strconv.Itoa(int(event.Time.Unix())),
			Data: b,
		})
	}

	c.emitter.Subscribe(eventKey, listener)

	for {
		select {
		case <-rw.CloseNotify:
			c.emitter.Unsubscribe(eventKey, listener)

			return
		}
	}
}
