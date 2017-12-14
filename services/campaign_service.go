// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package services

import (
	"github.com/asdine/storm"
	"github.com/euskadi31/cryptotrader/database/entity"
)

// CampaignServiceSave interface
type CampaignServiceSave interface {
	Save(data *entity.Campaign) error
}

// CampaignService struct
type CampaignService struct {
	db *storm.DB
}

// NewCampaignService constructor
func NewCampaignService(db *storm.DB) *CampaignService {
	return &CampaignService{
		db: db,
	}
}

// Save Campaign
func (s *CampaignService) Save(data *entity.Campaign) error {
	return s.db.Save(data)
}
