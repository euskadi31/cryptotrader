// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package database

import (
	"context"
	"errors"
	"net/http"

	"github.com/asdine/storm"
)

var (
	errContextIsNull     = errors.New("The context is null")
	errNotFountInContext = errors.New("The db is not found in context")
)

type key int

const (
	dbKey key = iota
)

func NewContext(ctx context.Context, db *storm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func NewHandler(db *storm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r != nil {
				r = r.WithContext(NewContext(r.Context(), db))
			}

			next.ServeHTTP(w, r)
		})
	}
}

// FromContext retruns db
func FromContext(ctx context.Context) (*storm.DB, error) {
	if ctx == nil {
		return nil, errContextIsNull
	}

	db, ok := ctx.Value(dbKey).(*storm.DB)
	if !ok {
		return nil, errNotFountInContext
	}

	return db, nil
}
