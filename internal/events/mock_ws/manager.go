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

	"github.com/gorilla/websocket"
	"github.com/twitchdev/twitch-cli/internal/events/mock_ws/mock_ws_server"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type ServerManager struct {
	serverList       *mock_ws_server.Servers
	reconnectTesting bool
	primaryServer    string
	port             int
	debugEnabled     bool
}

var serverManager *ServerManager

func StartWebsocketServers(enableDebug bool, port int) {
	serverManager = &ServerManager{
		serverList: &mock_ws_server.Servers{
			Servers: make(map[string]*mock_ws_server.WebSocketServer),
		},
		port:             port,
		reconnectTesting: false,
	}

	serverManager.debugEnabled = enableDebug

	// Start initial websocket server
	initialServer := &mock_ws_server.WebSocketServer{
		ServerId: util.RandomGUID()[:8],
		Status:   2,
		Clients: &mock_ws_server.Clients{
			Clients: make(map[string]*mock_ws_server.Client),
		},
		Upgrader:     websocket.Upgrader{},
		DebugEnabled: serverManager.debugEnabled,
	}
	serverManager.serverList.Put(initialServer.ServerId, initialServer)
	serverManager.primaryServer = initialServer.ServerId

	// Allow exit with Ctrl + C
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	m := http.NewServeMux()

	// Register URL handler
	m.HandleFunc("/ws", serverManager.wsPageHandler)

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

		// Spin up new server
		newServer := &mock_ws_server.WebSocketServer{
			ServerId: util.RandomGUID()[:8],
			Status:   2,
			Clients: &mock_ws_server.Clients{
				Clients: make(map[string]*mock_ws_server.Client),
			},
			Upgrader:     websocket.Upgrader{},
			DebugEnabled: serverManager.debugEnabled,
			// TODO: Include old clients and their subscriptions
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
	// TODO: Handle subscriptions
}
