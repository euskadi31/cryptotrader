// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"github.com/euskadi31/go-std"
)

// Trade struct
type Trade struct {
	Provider  string       `db:"provider" json:"provider"`
	ID        string       `db:"id" json:"id"`
	CreatedAt std.DateTime `db:"created_at" json:"created_at"`
	UpdatedAt std.DateTime `db:"updated_at" json:"updated_at"`
	DeletedAt std.DateTime `db:"deleted_at" json:"deleted_at"`
}
