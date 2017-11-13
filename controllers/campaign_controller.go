// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/euskadi31/cryptotrader/database/entity"
	"github.com/euskadi31/go-server"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// CampaignController struct
type CampaignController struct {
	db *storm.DB
}

// NewCampaignController constructor
func NewCampaignController(db *storm.DB) *CampaignController {
	return &CampaignController{
		db: db,
	}
}

// Mount implements server.Controller
func (c *CampaignController) Mount(r *server.Router) {
	r.AddRouteFunc("/api/v1/campaigns", c.GetListCampaignHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/api/v1/campaigns", c.PostCampaignHandler).Methods(http.MethodPost)
	r.AddRouteFunc("/api/v1/campaigns/{id:[0-9]+}", c.PutCampaignHandler).Methods(http.MethodPut)
	r.AddRouteFunc("/api/v1/campaigns/{id:[0-9]+}", c.GetCampaignHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/api/v1/campaigns/{id:[0-9]+}", c.DeleteCampaignHandler).Methods(http.MethodDelete)
}

// GetListCampaignHandler endpoint
func (c *CampaignController) GetListCampaignHandler(w http.ResponseWriter, r *http.Request) {
	var campaigns []*entity.Campaign

	if err := c.db.All(&campaigns); err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusOK, campaigns)
}

func (c *CampaignController) saveCampaign(r *http.Request) (*entity.Campaign, error) {
	campaign := &entity.Campaign{}

	id, isEdit := mux.Vars(r)["id"]

	if err := json.NewDecoder(r.Body).Decode(campaign); err != nil {
		return nil, err
	}

	if isEdit {
		i, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}

		campaign.ID = i
	}

	if err := c.db.Save(campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}

// PostCampaignHandler endpoint
func (c *CampaignController) PostCampaignHandler(w http.ResponseWriter, r *http.Request) {
	campaign, err := c.saveCampaign(r)
	if err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusCreated, campaign)
}

// PutCampaignHandler endpoint
func (c *CampaignController) PutCampaignHandler(w http.ResponseWriter, r *http.Request) {
	campaign, err := c.saveCampaign(r)
	if err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusOK, campaign)
}

func (c *CampaignController) fetchCampaign(r *http.Request) (*entity.Campaign, error) {
	id := mux.Vars(r)["id"]

	campaign := &entity.Campaign{}

	if err := c.db.One("ID", id, campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}

// GetCampaignHandler endpoint
func (c *CampaignController) GetCampaignHandler(w http.ResponseWriter, r *http.Request) {
	campaign, err := c.fetchCampaign(r)
	if err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusOK, campaign)
}

// DeleteCampaignHandler endpoint
func (c *CampaignController) DeleteCampaignHandler(w http.ResponseWriter, r *http.Request) {
	campaign, err := c.fetchCampaign(r)
	if err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	if err := c.db.DeleteStruct(campaign); err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusOK, campaign)
}
