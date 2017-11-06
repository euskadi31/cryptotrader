// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"github.com/euskadi31/go-std"
)

// User struct
type User struct {
	ID        int          `db:"id" json:"id"`
	Email     string       `db:"email" json:"email"`
	FirstName std.String   `db:"firstname" json:"firstname,omitempty"`
	LastName  std.String   `db:"lastname" json:"lastname,omitempty"`
	NickName  string       `db:"nickname" json:"nickname"`
	UserName  string       `db:"username" json:"username"`
	Password  string       `db:"password" json:"-"`
	Salt      string       `db:"salt" json:"salt"`
	IsEnabled bool         `db:"enabled" json:"enabled"`
	IsExpired bool         `db:"expired" json:"expired"`
	IsLocked  bool         `db:"locked" json:"locked"`
	Timezone  std.String   `db:"timezone" json:"timezone,omitempty"`
	Locale    std.String   `db:"locale" json:"locale,omitempty"`
	CreatedAt std.DateTime `db:"created_at" json:"created_at"`
	UpdatedAt std.DateTime `db:"updated_at" json:"updated_at"`
	DeletedAt std.DateTime `db:"deleted_at" json:"deleted_at"`
}
