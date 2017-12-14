// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package services

import (
	"github.com/asdine/storm"
	"github.com/euskadi31/cryptotrader/database/entity"
)

// OrderServiceSave interface
type OrderServiceSave interface {
	Save(data *entity.Order) error
}

// OrderService struct
type OrderService struct {
	db *storm.DB
}

// NewOrderService constructor
func NewOrderService(db *storm.DB) *OrderService {
	return &OrderService{
		db: db,
	}
}

// Save Order
func (s *OrderService) Save(data *entity.Order) error {
	return s.db.Save(data)
}
