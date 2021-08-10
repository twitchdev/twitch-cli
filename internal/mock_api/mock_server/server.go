// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints"
	"github.com/twitchdev/twitch-cli/internal/mock_api/generate"
	"github.com/twitchdev/twitch-cli/internal/mock_auth"
	"github.com/twitchdev/twitch-cli/internal/mock_units"
)

const MOCK_NAMESPACE = "/mock"
const UNITS_NAMESPACE = "/units"
const AUTH_NAMESPACE = "/auth"

func StartServer(port int) {
	m := http.NewServeMux()

	ctx := context.Background()

	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err.Error())
		return
	}

	firstTime := db.IsFirstRun()

	if firstTime {
		err := generate.Generate(25)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx = context.WithValue(ctx, "db", db)

	RegisterHandlers(m)
	s := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: m,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Print("Mock server started")
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-stop

	log.Print("shutting down ...\n")
	db.DB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func RegisterHandlers(m *http.ServeMux) {
	// all mock endpoints live in the /mock/ namespace
	for _, e := range endpoints.All() {
		// no auth requirements on this endpoint, so just add it manually
		if e.Path() == "/schedule/icalendar" {
			m.Handle(MOCK_NAMESPACE+e.Path(), loggerMiddleware(e))
			continue
		}
		m.Handle(MOCK_NAMESPACE+e.Path(), loggerMiddleware(authentication.AuthenticationMiddleware(e)))
	}
	for _, e := range mock_units.All() {
		m.Handle(UNITS_NAMESPACE+e.Path(), loggerMiddleware(e))
	}

	for _, e := range mock_auth.All() {
		m.Handle(AUTH_NAMESPACE+e.Path(), loggerMiddleware(e))
	}
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
