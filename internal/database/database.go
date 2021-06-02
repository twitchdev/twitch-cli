// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/util"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var dbFileName = "eventCache.db"

type CLIDatabase struct {
	DB *sqlx.DB
}

func NewConnection() (CLIDatabase, error) {
	db, err := getDatabase()
	if err != nil {
		return CLIDatabase{}, err
	}

	return CLIDatabase{DB: &db}, nil
}

func getDatabase() (sqlx.DB, error) {
	home, err := util.GetApplicationDir()
	if err != nil {
		return sqlx.DB{}, err
	}

	if viper.GetString("DB_FILENAME") != "" {
		dbFileName = viper.GetString("DB_FILENAME")
	}

	var path = filepath.Join(home, dbFileName)
	var needToInit = false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		needToInit = true
	}

	// open and force Foreign Key support ("fk=true")
	db, err := sqlx.Open("sqlite3", path+"?_fk=true")
	if err != nil {
		return sqlx.DB{}, err
	}

	if needToInit == true {
		err = initDatabase(*db)
		if err != nil {
			return sqlx.DB{}, err
		}
	}

	checkAndUpdate(*db)
	return *db, nil
}