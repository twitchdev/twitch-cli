// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints"
	"github.com/twitchdev/twitch-cli/internal/mock_api/generate"
	"github.com/twitchdev/twitch-cli/internal/mock_auth"
	"github.com/twitchdev/twitch-cli/internal/mock_units"
	"github.com/twitchdev/twitch-cli/internal/models"
)

const MOCK_NAMESPACE = "/mock"
const UNITS_NAMESPACE = "/units"
const AUTH_NAMESPACE = "/auth"

func StartServer(port int) error {
	m := http.NewServeMux()

	ctx := context.Background()

	db, err := database.NewConnection()
	if err != nil {
		return fmt.Errorf("Error connecting to database: %v", err.Error())
	}

	firstTime := db.IsFirstRun()

	if firstTime {
		err := generate.Generate(25)
		if err != nil {
			return err
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

	var serverErr error = nil

	go func() {
		log.Print("Mock server started")

		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				serverErr = err
				stop <- syscall.SIGINT // Simulate Ctrl+C
			}
		}
	}()

	<-stop

	if serverErr != nil {
		return serverErr
	}

	log.Print("shutting down ...\n")
	db.DB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		return err
	}

	return nil
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

	// For removed endpoints we don't have to worry about an actual handler, since its just gonna return 410 Gone
	for e := range endpoints.Gone() {
		m.Handle(MOCK_NAMESPACE+e, loggerMiddleware(nil))
	}
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.URL.Path)

		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		// Check for removed endpoints, which will return 410 Gone
		for goneEndpoint, methods := range endpoints.Gone() {
			if r.URL.Path == MOCK_NAMESPACE+goneEndpoint {
				validRemovedEndpoint := false
				for _, m := range methods {
					if strings.EqualFold(m, r.Method) {
						validRemovedEndpoint = true
					}
				}

				// In production, removed API URLs with no previously existing method return 404
				// e.g., "GET helix/tags/streams" returns 410, but "DELETE helix/tags/streams" returns 404
				if !validRemovedEndpoint {
					bytes, _ := json.Marshal(models.APIResponse{
						Error:   "Not Found",
						Status:  404,
						Message: "",
					})
					w.WriteHeader(http.StatusNotFound)
					w.Write(bytes)
				} else {
					bytes, _ := json.Marshal(models.APIResponse{
						Error:   "Gone",
						Status:  410,
						Message: "The API is deprecated.",
					})
					w.WriteHeader(410)
					w.Write(bytes)
				}
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
