// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"github.com/euskadi31/go-std"
)

// Campaign struct
type Campaign struct {
	ID          int          `db:"id,increment"`
	Provider    string       `db:"provider" json:"provider"`
	ProviderRef string       `db:"provider_ref" json:"provider_ref"`
	Volume      float64      `db:"volume" json:"volume"`
	CreatedAt   std.DateTime `db:"created_at" json:"created_at"`
	UpdatedAt   std.DateTime `db:"updated_at" json:"updated_at"`
	DeletedAt   std.DateTime `db:"deleted_at" json:"deleted_at"`
}
