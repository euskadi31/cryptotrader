// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"github.com/euskadi31/go-std"
)

// User struct
type User struct {
	ID        int          `storm:"id" json:"id"`
	Email     string       `storm:"email" json:"email"`
	FirstName std.String   `storm:"firstname" json:"firstname,omitempty"`
	LastName  std.String   `storm:"lastname" json:"lastname,omitempty"`
	NickName  string       `storm:"nickname" json:"nickname"`
	UserName  string       `storm:"username" json:"username"`
	Password  string       `storm:"password" json:"-"`
	Salt      string       `storm:"salt" json:"salt"`
	IsEnabled bool         `storm:"enabled" json:"enabled"`
	IsExpired bool         `storm:"expired" json:"expired"`
	IsLocked  bool         `storm:"locked" json:"locked"`
	Timezone  std.String   `storm:"timezone" json:"timezone,omitempty"`
	Locale    std.String   `storm:"locale" json:"locale,omitempty"`
	CreatedAt std.DateTime `storm:"created_at" json:"created_at"`
	UpdatedAt std.DateTime `storm:"updated_at" json:"updated_at"`
	DeletedAt std.DateTime `storm:"deleted_at" json:"deleted_at"`
}
