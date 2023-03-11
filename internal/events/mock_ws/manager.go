// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/twitchdev/twitch-cli/internal/events/mock_ws/mock_ws_server"
	"github.com/twitchdev/twitch-cli/internal/events/types"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type ServerManager struct {
	serverList       *util.List[mock_ws_server.WebSocketServer]
	reconnectTesting bool
	primaryServer    string
	port             int
	debugEnabled     bool
	strictMode       bool
}

var serverManager *ServerManager

func StartWebsocketServers(enableDebug bool, port int, strictMode bool) {
	serverManager = &ServerManager{
		serverList: &util.List[mock_ws_server.WebSocketServer]{
			Elements: make(map[string]*mock_ws_server.WebSocketServer),
		},
		port:             port,
		reconnectTesting: false,
		strictMode:       strictMode,
	}

	serverManager.debugEnabled = enableDebug

	// Start initial websocket server
	initialServer := &mock_ws_server.WebSocketServer{
		ServerId: util.RandomGUID()[:8],
		Status:   2,
		Clients: &util.List[mock_ws_server.Client]{
			Elements: make(map[string]*mock_ws_server.Client),
		},
		Upgrader:      websocket.Upgrader{},
		DebugEnabled:  serverManager.debugEnabled,
		Subscriptions: make(map[string][]mock_ws_server.Subscription),
		StrictMode:    serverManager.strictMode,
	}
	serverManager.serverList.Put(initialServer.ServerId, initialServer)
	serverManager.primaryServer = initialServer.ServerId

	// Allow exit with Ctrl + C
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	m := http.NewServeMux()

	// Register URL handler
	m.HandleFunc("/ws", serverManager.wsPageHandler)
	m.HandleFunc("/eventsub/subscriptions", serverManager.subscriptionPageHandler)

	// Start HTTP server
	go func() {
		// Listen to port
		listen, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
		if err != nil {
			log.Fatalf("Cannot start HTTP server: %v", err)
			return
		}

		log.Printf("Started WebSocket server on 127.0.0.1:%v", port)

		// Serve HTTP server
		if err := http.Serve(listen, m); err != nil {
			log.Fatalf("Cannot start HTTP server: %v", err)
			return
		}
	}()

	// Initalize RPC handler, to accept EventSub transports
	StartRPCListener()

	// TODO: Interactive shell maybe?

	<-stop // Wait for Ctrl + C
}

func (sm *ServerManager) HandleRPCServerCommand(messageBody string) bool {
	// Convert to struct for reading
	eventObj := models.EventsubResponse{}
	err := json.Unmarshal([]byte(messageBody), &eventObj)
	if err != nil {
		log.Printf("Error on RPC call (ServerCommand): Failed to parse command JSON: %v\nRaw: %v", err.Error(), messageBody)
		return false
	}

	if strings.EqualFold(eventObj.Subscription.Type, "mock.websocket.reconnect") {
		// Initiate reconnect testing
		log.Printf("Initiating reconnect testing...")

		if sm.reconnectTesting {
			log.Printf("Error on RPC call (ServerCommand): Cannot execute reconnect testing while its already in progress. Aborting.")
			return false
		}

		// Find current primary server
		originalPrimaryServer, ok := sm.serverList.Get(sm.primaryServer)
		if !ok {
			log.Printf("Error on RPC call (ServerCommand): Primary server not in server list.")
			return false
		}

		sm.reconnectTesting = true

		// Get the list of reconnect clients ready
		reconnectClients := originalPrimaryServer.GetCurrentSubscriptionsForReconnect()

		// Spin up new server
		newServer := &mock_ws_server.WebSocketServer{
			ServerId: util.RandomGUID()[:8],
			Status:   2,
			Clients: &util.List[mock_ws_server.Client]{
				Elements: make(map[string]*mock_ws_server.Client),
			},
			Upgrader:         websocket.Upgrader{},
			DebugEnabled:     serverManager.debugEnabled,
			Subscriptions:    make(map[string][]mock_ws_server.Subscription),
			StrictMode:       serverManager.strictMode,
			ReconnectClients: reconnectClients,
		}
		serverManager.serverList.Put(newServer.ServerId, newServer)

		// Switch manager's primary server to new one
		// Doing this before sending the reconnect messages emulates the Twitch's production load balancer, which will never send to servers shutting down.
		serverManager.primaryServer = newServer.ServerId

		// Notify primary server to restart (includes not accepting new clients)
		// This is in a goroutine so it doesn't hang the `twitch event trigger mock.websocket.reconnect` command
		go func() {
			originalPrimaryServer.InitiateRestart()

			// Remove server from server list
			serverManager.serverList.Delete(originalPrimaryServer.ServerId)

			if serverManager.debugEnabled {
				log.Printf(
					"Removed server [%v] from server list. New server list count: %v",
					originalPrimaryServer.ServerId,
					serverManager.serverList.Length(),
				)
			}

			serverManager.reconnectTesting = false

			log.Printf("Reconnect testing successful. Primary server is now [%v]\nYou may now execute reconnect testing again.", sm.primaryServer)
		}()

		return true
	} else {
		// Unknown command. Should technically never happen
		log.Printf("Error on RPC call (ServerCommand):  -- Unexpected server command: %v", eventObj.Subscription.Type)
		return false
	}
}

func (sm *ServerManager) wsPageHandler(w http.ResponseWriter, r *http.Request) {
	server, ok := sm.serverList.Get(sm.primaryServer)
	if !ok {
		log.Printf("Failed to find primary server [%v] when new client was accessing ws://127.0.0.1:%v/ws -- Aborting...", sm.primaryServer, sm.port)
		return
	}

	server.WsPageHandler(w, r)
}

func (sm ServerManager) subscriptionPageHandler(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)

	// OPTIONS method
	if method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Authorization, Client-Id, Twitch-Api-Token, X-Forwarded-Proto, X-Requested-With, X-Csrf-Token, Content-Type, X-Device-Id, X-Twitch-Vhscf, X-Forwarded-Ip")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "600")
		w.WriteHeader(http.StatusOK)
		return
	}

	// POST method
	if method == "POST" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("ratelimit-limit", "800")
		w.Header().Set("ratelimit-remaining", "799")
		w.Header().Set("ratelimit-reset", fmt.Sprintf("%d", time.Now().Unix()+1)) // 1 second from now

		var body mock_ws_server.SubscriptionPostRequest

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			handlerResponseErrorBadRequest(w, "error validating json")
			return
		}

		// Basic error checking
		if r.Header.Get("client-id") == "" {
			handlerResponseErrorUnauthorized(w, "Client-Id header required")
			return
		}
		if !strings.EqualFold(body.Transport.Method, "websocket") {
			handlerResponseErrorBadRequest(w, "The value specified in the 'method' field is not valid")
			return
		}
		if !sessionRegex.MatchString(body.Transport.SessionID) {
			handlerResponseErrorBadRequest(w, "The value specified in the 'session_id' field is not valid")
			return
		}
		if body.Type == "" {
			handlerResponseErrorBadRequest(w, "The value specified in the 'type' field is not valid")
			return
		}
		if body.Version == "" {
			handlerResponseErrorBadRequest(w, "The value specified in the 'version' field is not valid")
			return
		}

		// Check if the topic exists
		_, err = types.GetByTriggerAndTransport(body.Type, body.Transport.Method)
		if err != nil {
			handlerResponseErrorBadRequest(w, "The combination of values in the type and version fields is not valid")
			return
		}

		sessionRegexExec := sessionRegex.FindAllStringSubmatch(body.Transport.SessionID, -1)
		clientName := sessionRegexExec[0][2]

		// Get client and server
		server, ok := sm.serverList.Get(sessionRegexExec[0][1])
		if !ok {
			handlerResponseErrorBadRequest(w, "non-existent session_id")
			return
		}
		client, ok := server.Clients.Get(clientName)
		if !ok {
			handlerResponseErrorBadRequest(w, "non-existent session_id")
			return
		}

		server.MuSubscriptions.Lock()

		// Check for duplicate subscription
		for _, s := range server.Subscriptions[clientName] {
			if s.ClientID == r.Header.Get("client-id") && s.Type == body.Type && s.Version == body.Version {
				handlerResponseErrorConflict(w, "Subscription by the specified type and version combination for the specified Client ID already exists")
				server.MuSubscriptions.Unlock()
				return
			}
		}

		if len(server.Subscriptions[clientName]) >= 100 {
			handlerResponseErrorBadRequest(w, "You may only create 100 subscriptions within a single WebSocket connection")
			server.MuSubscriptions.Unlock()
			return
		}

		// Add subscription
		subscription := mock_ws_server.Subscription{
			SubscriptionID: util.RandomGUID(),
			ClientID:       r.Header.Get("client-id"),
			Type:           body.Type,
			Version:        body.Version,
			CreatedAt:      time.Now().UTC().Format(time.RFC3339Nano),
		}

		var subs []mock_ws_server.Subscription
		existingList, ok := server.Subscriptions[clientName]
		if ok {
			subs = existingList
		} else {
			subs = []mock_ws_server.Subscription{}
		}

		subs = append(subs, subscription)
		server.Subscriptions[clientName] = subs

		server.MuSubscriptions.Unlock()

		// Return 202 status code and response body
		w.WriteHeader(http.StatusAccepted)

		json.NewEncoder(w).Encode(&mock_ws_server.SubscriptionPostSuccessResponse{
			Body: mock_ws_server.SubscriptionPostSuccessResponseBody{
				ID:        subscription.SubscriptionID,
				Status:    "enabled",
				Type:      subscription.Type,
				Version:   subscription.Version,
				Condition: mock_ws_server.EmptyStruct{},
				CreatedAt: subscription.CreatedAt,
				Transport: mock_ws_server.SubscriptionTransport{
					Method:      "websocket",
					SessionID:   fmt.Sprintf("%v_%v", server.ServerId, clientName),
					ConnectedAt: client.ConnectedAtTimestamp,
				},
				Cost: 0,
			},
			Total:        0,
			MaxTotalCost: 10,
			TotalCost:    0,
		})

		if sm.debugEnabled {
			log.Printf(
				"Client ID [%v] created subscription [%v/%v] at subscription ID [%v]",
				r.Header.Get("client-id"),
				subscription.Type,
				subscription.Version,
				subscription.SubscriptionID,
			)
		}

		return
	}

	// DELETE method
	if method == "DELETE" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("ratelimit-limit", "800")
		w.Header().Set("ratelimit-remaining", "799")
		w.Header().Set("ratelimit-reset", fmt.Sprintf("%d", time.Now().Unix()+1)) // 1 second from now

		subscriptionId := r.URL.Query().Get("id")

		// Basic error checking
		if r.Header.Get("client-id") == "" {
			handlerResponseErrorUnauthorized(w, "Client-Id header required")
			return
		}
		if subscriptionId == "" {
			handlerResponseErrorBadRequest(w, "The id query parameter is required")
			return
		}

		server, ok := sm.serverList.Get(sm.primaryServer)
		if !ok {
			handlerResponseErrorInternalServerError(w, "Primary server not found in server list.")
			return
		}

		subFound := false

		server.MuSubscriptions.Lock()

		for client, clientSubscriptions := range server.Subscriptions {
			for i, subscription := range clientSubscriptions {
				if subscription.SubscriptionID == subscriptionId {
					subFound = true
					subsPart := make([]mock_ws_server.Subscription, 0)
					subsPart = append(subsPart, server.Subscriptions[client][:i]...)

					newSubs := append(subsPart, server.Subscriptions[client][i+1:]...)
					server.Subscriptions[client] = newSubs

					if sm.debugEnabled {
						log.Printf(
							"Deleted subscription [%v/%v] of ID [%v] owned by client ID [%v]",
							subscription.Type,
							subscription.Version,
							subscription.SubscriptionID,
							r.Header.Get("client-id"),
						)
					}
				}
			}
		}

		server.MuSubscriptions.Unlock()

		if subFound {
			// Return 204 status code
			w.WriteHeader(http.StatusNoContent)
		} else {
			// Return 404 not found
			w.WriteHeader(http.StatusNotFound)

			if sm.debugEnabled {
				log.Printf("Failed to delete subscription ID [%v] from client ID [%v]", subscriptionId, r.Header.Get("client-id"))
			}
		}

		return
	}

	// GET method
	if method == "GET" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("ratelimit-limit", "800")
		w.Header().Set("ratelimit-remaining", "799")
		w.Header().Set("ratelimit-reset", fmt.Sprintf("%d", time.Now().Unix()+1)) // 1 second from now

		// Basic error checking
		clientID := r.Header.Get("client-id")
		if clientID == "" {
			handlerResponseErrorUnauthorized(w, "Client-Id header required")
			return
		}

		server, ok := sm.serverList.Get(sm.primaryServer)
		if !ok {
			handlerResponseErrorInternalServerError(w, "Primary server not found in server list.")
			return
		}

		allSubscriptions := []mock_ws_server.SubscriptionPostSuccessResponseBody{}

		server.MuSubscriptions.Lock()

		for clientName, clientSubscriptions := range server.Subscriptions {
			for _, subscription := range clientSubscriptions {
				if subscription.ClientID == clientID {
					allSubscriptions = append(allSubscriptions, mock_ws_server.SubscriptionPostSuccessResponseBody{
						ID:        subscription.ClientID,
						Status:    "enabled",
						Type:      subscription.Type,
						Version:   subscription.Version,
						Condition: mock_ws_server.EmptyStruct{},
						CreatedAt: subscription.CreatedAt,
						Transport: mock_ws_server.SubscriptionTransport{
							Method:    "websocket",
							SessionID: fmt.Sprintf("%v_%v", server.ServerId, clientName),
						},
						Cost: 0,
					})
				}
			}
		}

		server.MuSubscriptions.Unlock()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&mock_ws_server.SubscriptionGetSuccessResponse{
			Total:        len(allSubscriptions),
			Data:         allSubscriptions,
			TotalCost:    0,
			MaxTotalCost: 10,
			Pagination:   mock_ws_server.EmptyStruct{},
		})

		return
	}

	// Fallback
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handlerResponseErrorBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	bytes, _ := json.Marshal(&mock_ws_server.SubscriptionPostErrorResponse{
		Error:   "Bad Request",
		Message: message,
		Status:  400,
	})
	w.Write(bytes)
}

func handlerResponseErrorUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	bytes, _ := json.Marshal(&mock_ws_server.SubscriptionPostErrorResponse{
		Error:   "Unauthorized",
		Message: message,
		Status:  401,
	})
	w.Write(bytes)
}

func handlerResponseErrorConflict(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusConflict)
	bytes, _ := json.Marshal(&mock_ws_server.SubscriptionPostErrorResponse{
		Error:   "Unauthorized",
		Message: message,
		Status:  409,
	})
	w.Write(bytes)
}

func handlerResponseErrorInternalServerError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	bytes, _ := json.Marshal(&mock_ws_server.SubscriptionPostErrorResponse{
		Error:   "Internal Server Error",
		Message: message,
		Status:  500,
	})
	w.Write(bytes)
}
