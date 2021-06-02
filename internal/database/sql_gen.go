// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

const SEP_AND = "and"
const SEP_OR = "or"

type InternalPagination struct {
	SQL string
	InternalCursor
	PaginationCursor string
}

type InternalCursor struct {
	Offset int `json:"o"`
	Limit  int `json:"l"`
}

// generates SELECT SQL for use with querying on an interface for easier querying. Generates the WHERE clause using a provided interface
func generateSQL(s string, i interface{}, seperator string) string {
	if seperator == "" {
		seperator = "and"
	}

	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	whereClause := []string{}
	for index := 0; index < t.NumField(); index++ {
		field := v.Field(index)
		dbField := ""
		db := t.Field(index).Tag.Get("db")
		if db == "" {
			continue
		}
		dbField = db

		// dbs or DB Select is used to ovverride the tag used when complex joins are used; allows for easier generation of SQL
		dbUpdate := t.Field(index).Tag.Get("dbs")
		if dbUpdate != "" {
			dbField = dbUpdate
		}

		switch field.Kind() {
		case reflect.String:
			if field.Interface().(string) == "" {
				continue
			}
			break
		case reflect.Ptr:
			if !field.Elem().IsValid() {
				continue
			}
			break
		case reflect.Bool:
			if field.Interface() == false {
				continue
			}
			break
		case reflect.Int:
			if field.Interface().(int) == 0 {
				continue
			}
			break
		default:
			break
		}

		whereClause = append(whereClause, fmt.Sprintf("%v = :%v", dbField, db))
	}

	if len(whereClause) == 0 {
		return s
	}

	w := strings.Join(whereClause, fmt.Sprintf(" %v ", seperator))
	s = fmt.Sprintf("%v where %v", s, w)
	return s
}

// Generates an Insert statement using a provided interface and optional upsert functionality
func generateInsertSQL(table string, pk string, i interface{}, upsert bool) string {
	t := reflect.TypeOf(i)

	insertClause := []string{}
	valuesClause := []string{}
	upsertClause := []string{}
	for index := 0; index < t.NumField(); index++ {
		field := t.Field(index)

		db := field.Tag.Get("db")
		if db == "" {
			continue
		}

		// dbi or Database Insert is used to omit specific fields from being inserted, such as during joins where certain fields come from other tables
		dbi := field.Tag.Get("dbi")
		if dbi == "false" {
			continue
		}

		insertClause = append(insertClause, fmt.Sprintf("%v", db))
		valuesClause = append(valuesClause, fmt.Sprintf(":%v", db))
		if upsert {
			upsertClause = append(upsertClause, fmt.Sprintf("%v=:%v", db, db))
		}
	}
	s := fmt.Sprintf("insert into %v (%v) values(%v)", table, strings.Join(insertClause, ", "), strings.Join(valuesClause, ", "))
	if upsert {
		s = fmt.Sprintf("%v on conflict(%v) do update set %v", s, pk, strings.Join(upsertClause, ", "))
	}

	return s
}

// Generates the respective pagination token
func generatePaginationSQLAndResponse(limit int, prev_cursor string, is_before bool) InternalPagination {
	ic := InternalCursor{}
	if limit == 0 {
		limit = 20
	}

	if prev_cursor != "" {
		b, err := base64.RawStdEncoding.DecodeString(prev_cursor)
		if err != nil {
			log.Print(err)
		}
		json.Unmarshal(b, &ic)
		if is_before {
			ic.Offset -= limit
		} else {
			ic.Offset += limit
		}
	}

	ic.Limit = limit

	if ic.Offset < 0 {
		return InternalPagination{}
	}

	b, _ := json.Marshal(ic)

	ip := InternalPagination{
		InternalCursor:   ic,
		PaginationCursor: base64.RawURLEncoding.EncodeToString(b),
		SQL:              fmt.Sprintf(" LIMIT %v OFFSET %v", ic.Limit, ic.Offset),
	}

	return ip
}