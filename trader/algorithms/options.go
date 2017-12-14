// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algorithms

import (
	"strings"
)

// Options type
type Options map[string]interface{}

// Merge Options
func (o Options) Merge(options Options) {
	for key, val := range options {
		o[key] = val
	}
}

// Set value
func (o Options) Set(key string, value interface{}) {
	o[strings.ToLower(key)] = value
}

// GetString value
func (o Options) GetString(key string) string {
	if o == nil {
		return ""
	}

	if value, ok := o[strings.ToLower(key)].(string); ok {
		return value
	}

	return ""
}

// GetBool value
func (o Options) GetBool(key string) bool {
	if o == nil {
		return false
	}

	if value, ok := o[strings.ToLower(key)].(bool); ok {
		return value
	}

	return false
}

// GetInt value
func (o Options) GetInt(key string) int {
	if o == nil {
		return 0
	}

	if value, ok := o[strings.ToLower(key)].(float64); ok {
		return int(value)
	} else if value, ok := o[strings.ToLower(key)].(int); ok {
		return value
	}

	return 0
}

// GetFloat value
func (o Options) GetFloat(key string) float64 {
	if o == nil {
		return 0
	}

	if value, ok := o[strings.ToLower(key)].(float64); ok {
		return value
	}

	return 0
}
