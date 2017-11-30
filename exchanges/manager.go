// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package exchanges

import (
	"errors"
)

var (
	ErrProviderNotFound = errors.New("provider not found")
	ErrManagerNotInit   = errors.New("manager is not init")
)

// Manager of ExchangeProvider
type Manager map[string]ExchangeProvider

// NewManager func
func NewManager() Manager {
	return Manager{}
}

// Get exchange by name
func (m Manager) Get(key string) (ExchangeProvider, error) {
	if m == nil {
		return nil, ErrManagerNotInit
	}

	exchange, ok := m[key]
	if ok == false {
		return nil, ErrProviderNotFound
	}

	return exchange, nil
}

// Add sets the key to value. It replaces any existing
// values.
func (m Manager) Add(exchange ExchangeProvider) {
	m[exchange.Name()] = exchange
}
