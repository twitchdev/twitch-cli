// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestGetDatabase(t *testing.T) {
	a := SetupTestEnv(t)
	p, _ := GetApplicationDir()

	dbFileName = viper.GetString("DB_FILENAME")

	// delete the existing temp db if it exists
	path := filepath.Join(p, dbFileName)
	err := os.Remove(path)

	// if the error is not that the file doesn't exist, fail the test
	if !os.IsNotExist(err) {
		a.Nil(err)
	}

	// since this creates a new db, will check those codepaths
	db, err := getDatabase()
	a.Nil(err)
	a.NotNil(db)

	// get again, making sure that this works
	db, err = getDatabase()
	a.Nil(err)
	a.NotNil(db)
}

func TestRetriveFromDB(t *testing.T) {
	a := SetupTestEnv(t)

	ecParams := *&EventCacheParameters{
		ID:        RandomGUID(),
		Event:     "foo",
		JSON:      "bar",
		FromUser:  "1234",
		ToUser:    "5678",
		Transport: "test",
		Timestamp: GetTimestamp().Format(time.RFC3339Nano),
	}

	err := InsertIntoDB(ecParams)
	a.Nil(err)

	dbResponse, err := GetEventByID(ecParams.ID)
	a.Nil(err)

	println(dbResponse.ID)
	a.NotNil(dbResponse)
	a.Equal("test", dbResponse.Transport)
}
