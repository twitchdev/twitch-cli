// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/twitchdev/twitch-cli/internal/events/types"
	rpc_handler "github.com/twitchdev/twitch-cli/internal/rpc"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type ServerManager struct {
	serverList       *util.List[WebSocketServer]
	reconnectTesting bool   // Indicates if the server is in the process of running a simulation server reconnect/restart
	primaryServer    string // The current primary server by its ID. This should be in serverList
	ip               string // IP the server will bind to
	port             int    // Port the server will bind to
	debugEnabled     bool   // Indicates if the server was started with --debug
	strictMode       bool   // Indicates if the server was started with --require-subscriptions
	sslEnabled       bool   // Indicates if the server was started with --ssl
	protocolHttp     string // String for the HTTP protocol URIs (http or https)
	protocolWs       string // String for the WS protocol URIs (ws or wss)
}

var serverManager *ServerManager

func StartWebsocketServer(enableDebug bool, ip string, port int, enableSSL bool, strictMode bool) {
	serverManager = &ServerManager{
		serverList: &util.List[WebSocketServer]{
			Elements: make(map[string]*WebSocketServer),
		},
		ip:               ip,
		port:             port,
		reconnectTesting: false,
		strictMode:       strictMode,
		sslEnabled:       enableSSL,
	}

	serverManager.debugEnabled = enableDebug

	// Start initial websocket server
	initialServer := &WebSocketServer{
		ServerId: util.RandomGUID()[:8],
		Status:   2,
		Clients: &util.List[Client]{
			Elements: make(map[string]*Client),
		},
		Upgrader:      websocket.Upgrader{},
		DebugEnabled:  serverManager.debugEnabled,
		Subscriptions: make(map[string][]Subscription),
		StrictMode:    serverManager.strictMode,

		ReconnectClients: &util.List[[]Subscription]{ // Empty and irrelevant at this point, but needed to avoid panic
			Elements: make(map[string]*[]Subscription),
		},
	}
	serverManager.serverList.Put(initialServer.ServerId, initialServer)
	serverManager.primaryServer = initialServer.ServerId

	// Allow exit with Ctrl + C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	m := http.NewServeMux()

	// Register URL handler
	m.HandleFunc("/ws", wsPageHandler)
	m.HandleFunc("/eventsub/subscriptions", subscriptionPageHandler)

	// Start HTTP server
	go func() {
		// Listen to port
		listen, err := net.Listen("tcp", fmt.Sprintf("%v:%v", ip, port))
		if err != nil {
			log.Fatalf("Cannot start HTTP server: %v", err)
			return
		}

		lightYellow := color.New(color.FgHiYellow).SprintFunc()
		lightRed := color.New(color.FgHiRed).SprintFunc()
		brightWhite := color.New(color.FgHiWhite).SprintFunc()

		// Serve HTTP server
		if serverManager.sslEnabled {
			serverManager.protocolHttp = "https"
			serverManager.protocolWs = "wss"

			home, err := util.GetApplicationDir()
			if err != nil {
				log.Fatalf("Cannot start HTTP server: %v", err)
				return
			}

			crtFile := filepath.Join(home, "localhost.crt")
			keyFile := filepath.Join(home, "localhost.key")
			_, crtFileErr := os.Stat(crtFile)
			_, keyFileErr := os.Stat(keyFile)
			if errors.Is(crtFileErr, os.ErrNotExist) || errors.Is(keyFileErr, os.ErrNotExist) {
				log.Fatalf(`%v
** Files must be placed in %v as %v and %v **
%v
** However, if you wish to generate the files using OpenSSL, run these commands: **
	openssl genrsa -out "%v" 2048
	openssl req -new -x509 -sha256 -key "%v" -out "%v" -days 365`,
					lightRed("ERROR: Missing one of the required SSL crt/key files."),
					brightWhite(home),
					brightWhite("localhost.crt"),
					brightWhite("localhost.key"),
					lightYellow("** Testing with Twitch CLI using SSL is meant for users experienced with SSL already, as these files must be added to your systems keychain to work without errors. **"),
					keyFile, keyFile, crtFile)
				return
			}

			printWelcomeMsg()

			if err := http.ServeTLS(listen, m, crtFile, keyFile); err != nil {
				log.Fatalf("Cannot start HTTP server: %v", err)
				return
			}
		} else {
			serverManager.protocolHttp = "http"
			serverManager.protocolWs = "ws"

			printWelcomeMsg()

			if err := http.Serve(listen, m); err != nil {
				log.Fatalf("Cannot start HTTP server: %v", err)
				return
			}
		}

	}()

	// Initalize RPC handler, to accept EventSub transports
	rpc := rpc_handler.RPCHandler{
		Port:     44747,
		Handlers: make(map[string]rpc_handler.HandlerCallback),
	}

	rpc.RegisterHandler("EventSubWebSocketReconnect", RPCReconnectHandler)
	rpc.RegisterHandler("EventSubWebSocketForwardEvent", RPCFireEventSubHandler)
	rpc.RegisterHandler("EventSubWebSocketCloseClient", RPCCloseHandler)
	rpc.RegisterHandler("EventSubWebSocketSubscription", RPCSubscriptionHandler)
	rpc.RegisterHandler("EventSubWebSocketKeepalive", RPCKeepaliveHandler)
	rpc.StartBackgroundServer()

	// TODO: Interactive shell maybe?

	<-stop // Wait for Ctrl + C
}

func printWelcomeMsg() {
	lightBlue := color.New(color.FgHiBlue).SprintFunc()
	lightGreen := color.New(color.FgHiGreen).SprintFunc()
	lightYellow := color.New(color.FgHiYellow).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	log.Printf(lightBlue("Started WebSocket server on %v:%v"), serverManager.ip, serverManager.port)
	if serverManager.strictMode {
		log.Println(lightBlue("--require-subscription enabled. Clients will have 10 seconds to subscribe before being disconnected."))
	}

	fmt.Println()

	log.Printf(yellow("Simulate subscribing to events at: %v://%v:%v/eventsub/subscriptions"), serverManager.protocolHttp, serverManager.ip, serverManager.port)
	log.Println(yellow("POST, GET, and DELETE are supported"))
	log.Println(yellow("For more info: https://dev.twitch.tv/docs/cli/websocket-event-command/#simulate-subscribing-to-mock-eventsub"))

	fmt.Println()

	log.Println(lightYellow("Events can be forwarded to this server from another terminal with --transport=websocket\nExample: \"twitch event trigger channel.ban --transport=websocket\""))
	fmt.Println()
	log.Println(lightYellow("You can send to a specific client after its connected with --session\nExample: \"twitch event trigger channel.ban --transport=websocket --session=e411cc1e_a2613d4e\""))

	fmt.Println()
	log.Println(lightGreen("For further usage information, please see our official documentation:\nhttps://dev.twitch.tv/docs/cli/websocket-event-command/"))
	fmt.Println()

	log.Printf(lightBlue("Connect to the WebSocket server at: ")+"%v://%v:%v/ws", serverManager.protocolWs, serverManager.ip, serverManager.port)
}

func wsPageHandler(w http.ResponseWriter, r *http.Request) {
	server, ok := serverManager.serverList.Get(serverManager.primaryServer)
	if !ok {
		log.Printf("Failed to find primary server [%v] when new client was accessing %v://%v:%v/ws -- Aborting...",
			serverManager.primaryServer, serverManager.protocolHttp, serverManager.ip, serverManager.port)
		return
	}

	server.WsPageHandler(w, r)
}

func subscriptionPageHandler(w http.ResponseWriter, r *http.Request) {
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

	// GET method
	if method == "GET" {
		subscriptionPageHandlerGet(w, r)
		return
	}

	// POST method
	if method == "POST" {
		subscriptionPageHandlerPost(w, r)
		return
	}

	// DELETE method
	if method == "DELETE" {
		subscriptionPageHandlerDelete(w, r)
		return
	}

	// Fallback
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func subscriptionPageHandlerGet(w http.ResponseWriter, r *http.Request) {
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

	server, ok := serverManager.serverList.Get(serverManager.primaryServer)
	if !ok {
		handlerResponseErrorInternalServerError(w, "Primary server not found in server list.")
		return
	}

	allSubscriptions := []SubscriptionPostSuccessResponseBody{}

	server.muSubscriptions.Lock()

	for clientName, clientSubscriptions := range server.Subscriptions {
		for _, subscription := range clientSubscriptions {
			disabledAndExpired := false // Production EventSub only shows disabled WebSocket subscriptions that were disabled under 1 hour ago
			if subscription.DisabledAt != nil && subscription.DisabledAt.Add(time.Hour).Before(util.GetTimestamp()) {
				disabledAndExpired = true
			}

			if clientID == "debug" || (subscription.ClientID == clientID && !disabledAndExpired) {
				allSubscriptions = append(allSubscriptions, SubscriptionPostSuccessResponseBody{
					ID:        subscription.SubscriptionID,
					Status:    subscription.Status,
					Type:      subscription.Type,
					Version:   subscription.Version,
					Condition: subscription.Conditions,
					CreatedAt: subscription.CreatedAt,
					Transport: SubscriptionTransport{
						Method:         "websocket",
						SessionID:      fmt.Sprintf("%v_%v", server.ServerId, clientName),
						ConnectedAt:    subscription.ClientConnectedAt,
						DisconnectedAt: subscription.ClientDisconnectedAt,
					},
					Cost: 0,
				})
			}
		}
	}

	server.muSubscriptions.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&SubscriptionGetSuccessResponse{
		Total:        len(allSubscriptions),
		Data:         allSubscriptions,
		TotalCost:    0,
		MaxTotalCost: 10,
		Pagination:   EmptyStruct{},
	})
}

func subscriptionPageHandlerPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("ratelimit-limit", "800")
	w.Header().Set("ratelimit-remaining", "799")
	w.Header().Set("ratelimit-reset", fmt.Sprintf("%d", time.Now().Unix()+1)) // 1 second from now

	var body SubscriptionPostRequest

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

	// Check if the topic was deprecated/removed
	for e, v := range types.RemovedEvents() {
		if body.Type == e && body.Version == v {
			handlerResponseErrorGone(w)
			return
		}
	}

	_, err = types.GetByTriggerAndTransportAndVersion(body.Type, body.Transport.Method, body.Version)
	if err != nil {
		handlerResponseErrorBadRequest(w, "The combination of values in the type and version fields is not valid")
		return
	}

	sessionRegexExec := sessionRegex.FindAllStringSubmatch(body.Transport.SessionID, -1)
	clientName := sessionRegexExec[0][2]

	// Get client and server
	server, ok := serverManager.serverList.Get(sessionRegexExec[0][1])
	if !ok {
		handlerResponseErrorBadRequest(w, "non-existent session_id")
		return
	}
	client, ok := server.Clients.Get(clientName)
	if !ok {
		handlerResponseErrorBadRequest(w, "non-existent session_id")
		return
	}

	server.muSubscriptions.Lock()

	// Check for duplicate subscription
	for _, s := range server.Subscriptions[clientName] {
		if s.ClientID == r.Header.Get("client-id") && s.Type == body.Type && s.Version == body.Version {
			handlerResponseErrorConflict(w, "Subscription by the specified type and version combination for the specified Client ID already exists")
			server.muSubscriptions.Unlock()
			return
		}
	}

	if len(server.Subscriptions[clientName]) >= 100 {
		handlerResponseErrorBadRequest(w, "You may only create 100 subscriptions within a single WebSocket connection")
		server.muSubscriptions.Unlock()
		return
	}

	// Add subscription
	subscription := Subscription{
		SubscriptionID:    util.RandomGUID(),
		ClientID:          r.Header.Get("client-id"),
		Type:              body.Type,
		Version:           body.Version,
		CreatedAt:         time.Now().UTC().Format(time.RFC3339Nano),
		Status:            STATUS_ENABLED, // https://dev.twitch.tv/docs/api/reference/#get-eventsub-subscriptions
		Conditions:        body.Condition,
		ClientConnectedAt: client.ConnectedAtTimestamp,
	}

	var subs []Subscription
	existingList, ok := server.Subscriptions[clientName]
	if ok {
		subs = existingList
	} else {
		subs = []Subscription{}
	}

	subs = append(subs, subscription)
	server.Subscriptions[clientName] = subs

	server.muSubscriptions.Unlock()

	// Return 202 status code and response body
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(&SubscriptionPostSuccessResponse{
		Data: []SubscriptionPostSuccessResponseBody{
			{
				ID:        subscription.SubscriptionID,
				Status:    subscription.Status,
				Type:      subscription.Type,
				Version:   subscription.Version,
				Condition: subscription.Conditions,
				CreatedAt: subscription.CreatedAt,
				Transport: SubscriptionTransport{
					Method:      "websocket",
					SessionID:   fmt.Sprintf("%v_%v", server.ServerId, clientName),
					ConnectedAt: client.ConnectedAtTimestamp,
				},
				Cost: 0,
			},
		},
		Total:        0,
		MaxTotalCost: 10,
		TotalCost:    0,
	})

	if serverManager.debugEnabled {
		log.Printf(
			"Client ID [%v] created subscription [%v/%v] at subscription ID [%v]",
			r.Header.Get("client-id"),
			subscription.Type,
			subscription.Version,
			subscription.SubscriptionID,
		)
	}
}

func subscriptionPageHandlerDelete(w http.ResponseWriter, r *http.Request) {
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

	server, ok := serverManager.serverList.Get(serverManager.primaryServer)
	if !ok {
		handlerResponseErrorInternalServerError(w, "Primary server not found in server list.")
		return
	}

	subFound := false

	server.muSubscriptions.Lock()

	for client, clientSubscriptions := range server.Subscriptions {
		for i, subscription := range clientSubscriptions {
			if subscription.SubscriptionID == subscriptionId {
				subFound = true
				subsPart := make([]Subscription, 0)
				subsPart = append(subsPart, server.Subscriptions[client][:i]...)

				newSubs := append(subsPart, server.Subscriptions[client][i+1:]...)
				server.Subscriptions[client] = newSubs

				if serverManager.debugEnabled {
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

	server.muSubscriptions.Unlock()

	if subFound {
		// Return 204 status code
		w.WriteHeader(http.StatusNoContent)
	} else {
		// Return 404 not found
		w.WriteHeader(http.StatusNotFound)

		if serverManager.debugEnabled {
			log.Printf("Failed to delete subscription ID [%v] from client ID [%v]", subscriptionId, r.Header.Get("client-id"))
		}
	}
}

func handlerResponseErrorBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	bytes, _ := json.Marshal(&SubscriptionPostErrorResponse{
		Error:   "Bad Request",
		Message: message,
		Status:  400,
	})
	w.Write(bytes)
}

func handlerResponseErrorUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	bytes, _ := json.Marshal(&SubscriptionPostErrorResponse{
		Error:   "Unauthorized",
		Message: message,
		Status:  401,
	})
	w.Write(bytes)
}

func handlerResponseErrorConflict(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusConflict)
	bytes, _ := json.Marshal(&SubscriptionPostErrorResponse{
		Error:   "Conflict",
		Message: message,
		Status:  409,
	})
	w.Write(bytes)
}

func handlerResponseErrorInternalServerError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	bytes, _ := json.Marshal(&SubscriptionPostErrorResponse{
		Error:   "Internal Server Error",
		Message: message,
		Status:  500,
	})
	w.Write(bytes)
}

func handlerResponseErrorGone(w http.ResponseWriter) {
	w.WriteHeader(http.StatusGone)
	bytes, _ := json.Marshal(&SubscriptionPostErrorResponse{
		Error:   "Gone",
		Message: "This subscription type is not available.",
		Status:  410,
	})
	w.Write(bytes)
}
