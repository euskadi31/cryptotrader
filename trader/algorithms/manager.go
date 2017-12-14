// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algorithms

import (
	"errors"
)

// Errors
var (
	ErrAlgorithmNotFound = errors.New("algorithm not found")
	ErrManagerNotInit    = errors.New("manager is not init")
)

// Manager of Algorithm
type Manager map[string]Algorithm

// NewManager func
func NewManager() Manager {
	return Manager{}
}

// Get algorithm by name
func (m Manager) Get(key string) (Algorithm, error) {
	if m == nil {
		return nil, ErrManagerNotInit
	}

	algorithm, ok := m[key]
	if ok == false {
		return nil, ErrAlgorithmNotFound
	}

	return algorithm, nil
}

// Has Algorithm
func (m Manager) Has(key string) bool {
	_, ok := m[key]

	return ok
}

// Add sets the key to value. It replaces any existing
// values.
func (m Manager) Add(algorithm Algorithm) {
	m[algorithm.Name()] = algorithm
}
