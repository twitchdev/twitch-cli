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

func generateSQL(s string, i interface{}, seperator string) string {
	if seperator == "" {
		seperator = "and"
	}

	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	whereClause := []string{}
	for index := 0; index < t.NumField(); index++ {
		field := v.Field(index)

		db := t.Field(index).Tag.Get("db")
		if db == "" {
			continue
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
		default:
			println(v.Field(index).Kind())
		}

		whereClause = append(whereClause, fmt.Sprintf("%v = :%v", db, db))
	}
	w := strings.Join(whereClause, fmt.Sprintf(" %v ", seperator))
	s = fmt.Sprintf("%v where %v", s, w)
	return s
}

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
