// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type Query struct {
	Limit  int
	Cursor string
	InternalPagination
	DB *sqlx.DB
}

// NewQuery handles the logic for generating the pagination token to pass alongside the DB queries for easier access
func (c CLIDatabase) NewQuery(r *http.Request, max_limit int) *Query {
	return c.NewQueryWithDefaultLimit(r, max_limit, 20)
}

func (c CLIDatabase) NewQueryWithDefaultLimit(r *http.Request, max_limit int, default_limit int) *Query {
	p := Query{DB: c.DB}
	if r == nil {
		return &p
	}

	ic := InternalCursor{}

	query := r.URL.Query()
	a := query.Get("after")
	f := query.Get("first")
	b := query.Get("before")

	isBefore := false
	if b != "" {
		isBefore = true
	}

	if len(a) > 0 {
		p.Cursor = a
	}

	first, _ := strconv.Atoi(f)
	if first > max_limit || first <= 0 {
		first = default_limit
	}
	p.Limit = int(first)

	if a != "" {
		b, err := base64.RawStdEncoding.DecodeString(a)
		if err != nil {
			return &p
		}
		err = json.Unmarshal(b, &ic)
		if err != nil {
			return &p
		}

		if isBefore {
			ic.Offset -= first
		} else {
			ic.Offset += first
		}
	}

	ic.Limit = first

	if ic.Offset < 0 {
		return &p
	}

	body, _ := json.Marshal(ic)

	ip := InternalPagination{
		InternalCursor:   ic,
		PaginationCursor: base64.RawURLEncoding.EncodeToString(body),
		SQL:              fmt.Sprintf(" LIMIT %v OFFSET %v", ic.Limit, ic.Offset),
	}

	p.InternalPagination = ip
	return &p
}
