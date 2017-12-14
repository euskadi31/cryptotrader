// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/euskadi31/cryptotrader/trader/algorithms"
	"github.com/euskadi31/go-server"
)

// AlgorithmController struct
type AlgorithmController struct {
	algorithms algorithms.Manager
}

// NewAlgorithmController constructor
func NewAlgorithmController(algorithms algorithms.Manager) *AlgorithmController {
	return &AlgorithmController{
		algorithms: algorithms,
	}
}

// Mount implements server.Controller
func (c *AlgorithmController) Mount(r *server.Router) {
	r.AddRouteFunc("/api/v1/algorithms", c.GetListHandler).Methods(http.MethodGet)
}

// GetListHandler endpoint
func (c *AlgorithmController) GetListHandler(w http.ResponseWriter, r *http.Request) {
	server.JSON(w, http.StatusOK, c.algorithms)
}
