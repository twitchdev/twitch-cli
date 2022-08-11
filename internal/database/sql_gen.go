// Copyright Amazon.com, Inc. or its affiliates. All Rights Reseved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
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
func generateSQL(s string, i interface{}, separator string) string {
	if separator == "" {
		separator = "and"
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

		switch f := field.Interface().(type) {
		case sql.NullString:
			if !f.Valid {
				continue
			}
		case string:
			if field.Interface().(string) == "" {
				continue
			}
		case bool:
			if !f {
				continue
			}
		case int:
			if f == 0 {
				continue
			}
		case float64:
			if f == 0.0 {
				continue
			}
		case time.Time:
			if f.IsZero() {
				continue
			}
		default:
			break
		}

		if v.Field(index).Kind() == reflect.Ptr {
			if v.Field(index).IsNil() {
				continue
			}
		}
		whereClause = append(whereClause, fmt.Sprintf("%v = :%v", dbField, db))
	}

	if len(whereClause) == 0 {
		return s
	}

	w := strings.Join(whereClause, fmt.Sprintf(" %v ", separator))
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

func generateUpdateSQL(table string, pk []string, i interface{}) string {
	updateClause := generateStructUpdateString(i)
	s := fmt.Sprintf("update %v set %v", table, strings.Join(updateClause, ", "))
	if len(pk) == 0 {
		return s
	}

	whereClause := []string{}
	for _, key := range pk {
		whereClause = append(whereClause, fmt.Sprintf("%v=:%v", key, key))
	}
	return fmt.Sprintf("%v where %v", s, strings.Join(whereClause, " AND "))
}

func generateStructUpdateString(i interface{}) []string {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	updateClause := []string{}
	for index := 0; index < t.NumField(); index++ {
		field := v.Field(index)

		switch v.Field(index).Kind() {
		case reflect.Ptr:
			if v.Field(index).IsNil() {
				continue
			}
		case reflect.Struct:
			nested := generateStructUpdateString(field.Interface())
			updateClause = append(updateClause, nested...)
		}

		db := t.Field(index).Tag.Get("db")
		if db == "" {
			continue
		}

		// dbi or Database Insert is used to omit specific fields from being inserted, such as during joins where certain fields come from other tables
		dbi := t.Field(index).Tag.Get("dbi")
		if dbi == "false" {
			continue
		}

		// if you need to force an update to a null value, force will do so
		if dbi == "force" {
			updateClause = append(updateClause, fmt.Sprintf("%v=:%v", db, db))
			continue
		}

		switch f := field.Interface().(type) {
		case sql.NullString:
			if !f.Valid {
				continue
			}
		case string:
			if field.Interface().(string) == "" {
				continue
			}
		case bool:
			if !f {
				continue
			}
		case int:
			if f == 0 {
				continue
			}
		default:
			break
		}

		updateClause = append(updateClause, fmt.Sprintf("%v=:%v", db, db))
	}
	return updateClause
}
