// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"log"
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

	for i := 0; i <= 5; i++ {
		// open and force Foreign Key support ("fk=true")
		db, err := sqlx.Open("sqlite3", path+"?_fk=true&cache=shared")
		if err != nil {
			log.Print(i)
			if i == 5 {
				return sqlx.DB{}, err
			}
			continue
		}

		if needToInit {
			err = initDatabase(*db)
			if err != nil {
				log.Printf("%#v", err)
				return sqlx.DB{}, err
			}
		}
		db.SetMaxOpenConns(1)
		checkAndUpdate(*db)
		return *db, nil
	}
	return sqlx.DB{}, nil
}
