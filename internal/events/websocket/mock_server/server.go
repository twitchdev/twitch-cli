package mock_server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

// Minimum time between messages before the server disconnects a client.
const KEEPALIVE_TIMEOUT_SECONDS = 10

type WebSocketServer struct {
	ServerId string // Int representing the ID of the server
	//ConnectionUrl string // Server's url for people to connect to. Used for messaging in reconnect testing
	DebugEnabled bool // Display debug messages; --debug
	StrictMode   bool // Force stricter production-like qualities; --strict
	Upgrader     websocket.Upgrader

	Clients   *util.List[Client] // All connected clients
	muClients sync.Mutex         // Mutex for WebSocketServer.Clients

	Status   int        // 0 = shut down; 1 = shutting down; 2 = online
	muStatus sync.Mutex // Mutex for WebSocketServer.Status

	Subscriptions   map[string][]Subscription // Active subscriptions on this server -- Accessed via Subscriptions[clientName]
	muSubscriptions sync.Mutex                // Mutex for WebSocketServer.Subscriptions

	ReconnectClients   *util.List[[]Subscription] // Clients that were part of the last server
	muReconnectClients sync.Mutex                 // Mutex for WebSocketServer.ReconnectClients
}

func (ws *WebSocketServer) WsPageHandler(w http.ResponseWriter, r *http.Request) {
	// This next line is required to disable CORS checking. No sense in caring in a test environment.
	ws.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("[[websocket upgrade err]] ", err)
		return
	}
	defer conn.Close()

	// Connection successful. WebSocket is open.

	// Get connected at time and set automatic read timeout
	connectedAtTimestamp := time.Now().UTC().Format(time.RFC3339Nano)
	conn.SetReadDeadline(time.Now().Add(time.Second * KEEPALIVE_TIMEOUT_SECONDS))

	client := &Client{
		clientName:           util.RandomGUID()[:8],
		conn:                 conn,
		ConnectedAtTimestamp: connectedAtTimestamp,
		connectionUrl:        fmt.Sprintf("%v://%v/ws", serverManager.protocolHttp, r.Host),
		keepAliveChanOpen:    false,
		pingChanOpen:         false,
	}

	if r.URL.Query().Get("reconnect_id") != "" {
		reconnectIdBytes, err := base64.StdEncoding.DecodeString(r.URL.Query().Get("reconnect_id") + "=")
		if err != nil {
			if ws.DebugEnabled {
				log.Printf("Could not decode base64 reconnect_id query parameter: '%v'", r.URL.Query().Get("reconnect_id"))
			}
		} else {
			reconnectId := string(reconnectIdBytes)

			ws.muReconnectClients.Lock()

			subscriptions, ok := ws.ReconnectClients.Get(reconnectId)
			if ok { // User had subscriptions carry over
				ws.Subscriptions[client.clientName] = *subscriptions
			}

			ws.ReconnectClients.Delete(reconnectId)

			if ws.DebugEnabled {
				log.Printf("Reconnected client [%v] was assigned %v subscriptions", client.clientName, len(ws.Subscriptions[client.clientName]))
			}

			ws.muReconnectClients.Unlock()
		}
	}

	// Disconnect the user if the server is in reconnect phase
	ws.muStatus.Lock()
	if ws.Status != 2 {
		// This is the closest we can get to the production environment, as there's no way to route people to a shutting down server
		log.Printf("New client trying to connect while websocket server in reconnect phase. Disconnecting them.")
		client.CloseDirty()
		// No handleClientConnectionClose because client is not in clients list, and chan loop was not set up yet.

		ws.muStatus.Unlock()
		return
	}

	// TODO: Check if user is connected to the old server, and if they are then disconnect them from old server with close frame 4004

	ws.muClients.Lock()
	// Add to the client connections list
	ws.Clients.Put(client.clientName, client)
	ws.muClients.Unlock()

	// This is put after ws.Clients.Put to make sure the client gets included in the list before InitiateRestart() kicks everyone out
	// Avoids any possible rare edge cases. This ain't production but I can still be safe :)
	ws.muStatus.Unlock()

	log.Printf("Client connected [%v]", client.clientName)
	ws.printConnections()

	// Send welcome message
	welcomeMsg, _ := json.Marshal(
		WelcomeMessage{
			Metadata: MessageMetadata{
				MessageID:        util.RandomGUID(),
				MessageType:      "session_welcome",
				MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
			},
			Payload: WelcomeMessagePayload{
				Session: WelcomeMessagePayloadSession{
					ID:                      fmt.Sprintf("%v_%v", ws.ServerId, client.clientName),
					Status:                  "connected",
					KeepaliveTimeoutSeconds: KEEPALIVE_TIMEOUT_SECONDS,
					ReconnectUrl:            nil,
					ConnectedAt:             connectedAtTimestamp,
				},
			},
		},
	)
	client.SendMessage(websocket.TextMessage, welcomeMsg)

	// Check if any subscriptions are sent after 10 seconds.
	// Strict mode only
	client.mustSubscribeTimer = time.NewTimer(10 * time.Second)
	if ws.StrictMode {
		go func() {
			select {
			case <-client.mustSubscribeTimer.C:
				if len(ws.Subscriptions[client.clientName]) == 0 {
					client.CloseWithReason(closeConnectionUnused)
					ws.handleClientConnectionClose(client, closeConnectionUnused)

					return
				}
			}
		}()
	}

	// Set up ping/pong and keepalive handling
	client.keepAliveTimer = time.NewTicker(10 * time.Second)
	client.pingTimer = time.NewTicker(5 * time.Second)
	client.keepAliveLoopChan = make(chan struct{})
	client.pingLoopChan = make(chan struct{})
	client.keepAliveChanOpen = true
	client.pingChanOpen = true

	// Set pong handler. Resets the read deadline when pong is received.
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(time.Second * KEEPALIVE_TIMEOUT_SECONDS))
		return nil
	})

	// Keep Alive message loop
	go func() {
		for {
			select {
			case <-client.keepAliveLoopChan:
				client.keepAliveTimer.Stop()
				client.keepAliveLoopChan = nil
				return

			case <-client.keepAliveTimer.C: // Send KeepAlive message
				keepAliveMsg, _ := json.Marshal(
					KeepaliveMessage{
						Metadata: MessageMetadata{
							MessageID:        util.RandomGUID(),
							MessageType:      "session_keepalive",
							MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
						},
						Payload: KeepaliveMessagePayload{},
					},
				)
				err := client.SendMessage(websocket.TextMessage, keepAliveMsg)
				if err != nil {
					client.CloseWithReason(closeNetworkError)
				}

				if ws.DebugEnabled {
					log.Printf("Sent session_keepalive to client [%s]", client.clientName)
				}
			}
		}
	}()

	// Ping/pong handler loop
	go func() {
		for {
			select {
			case <-client.pingLoopChan:
				client.pingTimer.Stop()
				client.pingLoopChan = nil
				return

			case <-client.pingTimer.C: // Send ping
				err := client.SendMessage(websocket.PingMessage, []byte{})
				if err != nil {
					ws.muClients.Lock()
					client.CloseWithReason(closeClientFailedPingPong)
					ws.handleClientConnectionClose(client, closeClientFailedPingPong)
					ws.muClients.Unlock()
				}

				if ws.DebugEnabled {
					log.Printf("Sent pong to client [%s]", client.clientName)
				}

			}
		}
	}()

	// Wait for message
	for {
		// Reset timeout upon every message, no matter what it is.
		client.conn.SetReadDeadline(time.Now().Add(time.Second * KEEPALIVE_TIMEOUT_SECONDS))

		mt, message, err := conn.ReadMessage()
		if err != nil && ws.Status != 0 { // If server is shut down, clients should already be disconnectd.
			log.Printf("read err [%v]: %v", client.clientName, err)

			ws.muClients.Lock()
			client.CloseWithReason(closeClientDisconnected)
			ws.handleClientConnectionClose(client, closeClientDisconnected)
			ws.muClients.Unlock()
			break
		}

		if ws.Status == 2 { // Only care about this when the server is running
			log.Printf("Disconnecting client [%v] due to received inbound traffic.\nMessage[%d]: %s", client.clientName, mt, message)

			ws.muClients.Lock()
			client.CloseWithReason(closeClientSentInboundTraffic)
			ws.handleClientConnectionClose(client, closeClientSentInboundTraffic)
			ws.muClients.Unlock()
		}

		break
	}
}

// Gets client subscriptions to be transfered to another server. Used during reconnect testing.
func (ws *WebSocketServer) GetCurrentSubscriptionsForReconnect() *util.List[[]Subscription] {
	reconnectClients := &util.List[[]Subscription]{
		Elements: make(map[string]*[]Subscription),
	}

	ws.muSubscriptions.Lock()

	for clientName, clientSubscriptions := range ws.Subscriptions {
		for _, subscription := range clientSubscriptions {
			reconnectReference := fmt.Sprintf("%v_%v", ws.ServerId, clientName)

			oldReconnectSubs, ok := reconnectClients.Get(reconnectReference)
			if !ok {
				oldReconnectSubs = &[]Subscription{}
			}

			// Add to oldReconnectSubs
			*oldReconnectSubs = append(*oldReconnectSubs, subscription)

			// Return to list
			reconnectClients.Put(reconnectReference, oldReconnectSubs)
		}
	}

	ws.muSubscriptions.Unlock()

	return reconnectClients
}

func (ws *WebSocketServer) InitiateRestart() {
	// Set status to shutting down; Stop accepting new clients
	ws.muStatus.Lock()
	ws.Status = 1
	ws.muStatus.Unlock()

	ws.muClients.Lock()

	if ws.DebugEnabled {
		log.Printf("Sending reconnect notices to [%v] clients", ws.Clients.Length())
	}

	// Send reconnect messages and disable timers on all clients
	for _, client := range ws.Clients.All() {
		// Disable keepalive and subscription timers
		close(client.keepAliveLoopChan)
		client.keepAliveChanOpen = false
		client.mustSubscribeTimer.Stop()

		// Send reconnect notice
		sessionId := fmt.Sprintf("%v_%v", ws.ServerId, client.clientName)
		reconnectId := base64.StdEncoding.EncodeToString([]byte(sessionId))
		reconnectId = reconnectId[:len(reconnectId)-1]
		clientConnectionUrl := strings.Replace(client.connectionUrl, "http://", "ws://", -1)
		clientConnectionUrl = strings.Replace(clientConnectionUrl, "https://", "wss://", -1)
		reconnectMsg, _ := json.Marshal(
			ReconnectMessage{
				Metadata: MessageMetadata{
					MessageID:        util.RandomGUID(),
					MessageType:      "session_reconnect",
					MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
				},
				Payload: ReconnectMessagePayload{
					Session: ReconnectMessagePayloadSession{
						ID:                      sessionId,
						Status:                  "reconnecting",
						KeepaliveTimeoutSeconds: nil,
						ReconnectUrl:            fmt.Sprintf("%v?reconnect_id=%v", clientConnectionUrl, reconnectId),
						ConnectedAt:             client.ConnectedAtTimestamp,
					},
				},
			},
		)

		err := client.SendMessage(websocket.TextMessage, reconnectMsg)
		if err != nil {
			log.Printf("Error building session_reconnect JSON for client [%v]: %v", client.clientName, err.Error())
		}
	}

	log.Printf("Reconnect notices sent for server [%v]", ws.ServerId)
	log.Printf("Will disconnect all existing clients in 30 seconds...")

	ws.muClients.Unlock()

	// Wait 30 seconds
	time.Sleep(30 * time.Second)

	// Change server status to 0
	// This is done before disconnects because the read loop will err out due to the close message, which gets printed unless this is zero.
	ws.Status = 0

	// Disconnect everyone with reconnect close message
	for _, client := range ws.Clients.All() {
		ws.muClients.Lock()
		client.CloseWithReason(closeReconnectGraceTimeExpired)
		ws.handleClientConnectionClose(client, closeReconnectGraceTimeExpired)
		ws.muClients.Unlock()
	}

	log.Printf("All users disconnected from server [%v]", ws.ServerId)
}

func (ws *WebSocketServer) HandleRPCEventSubForwarding(eventsubBody string, clientName string) (bool, string) {
	// If --session is used, make sure the client exists
	if clientName != "" {
		_, ok := ws.Clients.Get(strings.ToLower(clientName))
		if !ok {
			msg := fmt.Sprintf("Error executing remote triggered EventSub: Client [%v] does not exist on server [%v]", clientName, ws.ServerId)
			log.Println(msg)
			return false, msg
		}
	}

	if ws.Clients.Length() == 0 {
		msg := fmt.Sprintf("Warning for remote triggered EventSub: No clients in server [%v]", ws.ServerId)
		log.Println(msg)
		return false, msg
	}

	// Convert to struct for editing
	eventObj := models.EventsubResponse{}
	err := json.Unmarshal([]byte(eventsubBody), &eventObj)
	if err != nil {
		msg := fmt.Sprintf("Error reading JSON forwarded from EventSub: %v\nRaw: %v", err.Error(), eventsubBody)
		log.Println(msg)
		return false, msg
	}

	didSend := false

	for _, client := range ws.Clients.All() {
		if clientName != "" && !strings.EqualFold(strings.ToLower(clientName), clientName) {
			// When --session is used, only send to that client
			continue
		}

		// If this is a Revocation message (user.authorization.revoke), set it as revoked
		if eventObj.Subscription.Type == "user.authorization.revoke" {
			if serverManager.debugEnabled {
				log.Printf("Attempting to revoke subscription [%v]", eventObj.Subscription.ID)
			}

			ws.muSubscriptions.Lock()
			foundClientId := ""
			for client, clientSubscriptions := range ws.Subscriptions {
				if foundClientId != "" {
					break
				}

				for i, sub := range clientSubscriptions {
					if sub.SubscriptionID == eventObj.Subscription.ID {
						foundClientId = sub.ClientID

						ws.Subscriptions[client][i].Status = STATUS_AUTHORIZATION_REVOKED
						break
					}
				}
			}
			ws.muSubscriptions.Unlock()

			if foundClientId != "" {
				log.Printf("Subscription ID [%v], belonging to Client ID [%v], has been revoked.", eventObj.Subscription.ID, foundClientId)
			} else {
				msg := fmt.Sprintf("Failed to revoke Subscription ID [%v]: Subscription by that ID does not exist.", eventObj.Subscription.ID)
				log.Println(msg)
				return false, msg
			}
		}

		// Check for subscriptions when running with --require-subscription
		if ws.StrictMode {
			found := false
			for _, clientSubscriptions := range ws.Subscriptions {
				if found {
					break
				}

				for _, sub := range clientSubscriptions {
					if sub.SessionClientName == client.clientName && sub.Type == eventObj.Subscription.Type && sub.Version == eventObj.Subscription.Version {
						found = true
					}
				}
			}

			if !found {
				continue
			}
		}

		// Change payload's subscription.transport.session_id to contain the correct Session ID
		eventObj.Subscription.Transport.SessionID = fmt.Sprintf("%v_%v", ws.ServerId, client.clientName)

		// Build notification message
		notificationMsg, err := json.Marshal(
			NotificationMessage{
				Metadata: MessageMetadata{
					MessageID:           util.RandomGUID(),
					MessageType:         "notification",
					MessageTimestamp:    time.Now().UTC().Format(time.RFC3339Nano),
					SubscriptionType:    eventObj.Subscription.Type,
					SubscriptionVersion: eventObj.Subscription.Version,
				},
				Payload: eventObj,
			},
		)
		if err != nil {
			msg := fmt.Sprintf("Error building JSON for client [%v]: %v", client.clientName, err.Error())
			log.Println(msg)
			return false, msg
		}

		client.SendMessage(websocket.TextMessage, notificationMsg)
		log.Printf("Sent [%v / %v] to client [%v]", eventObj.Subscription.Type, eventObj.Subscription.Version, client.clientName)

		didSend = true
	}

	if !didSend {
		msg := fmt.Sprintf("Error executing remote triggered EventSub: No clients with the subscribed to [%v / %v]", eventObj.Subscription.Type, eventObj.Subscription.Version)
		log.Println(msg)
		return false, msg
	}

	return true, ""
}

func (ws *WebSocketServer) handleClientConnectionClose(client *Client, closeReason *CloseMessage) {
	// Prevent further looping
	client.mustSubscribeTimer.Stop()
	if client.keepAliveChanOpen {
		close(client.keepAliveLoopChan)
		client.keepAliveChanOpen = false
	}
	if client.pingChanOpen {
		close(client.pingLoopChan)
		client.pingChanOpen = false
	}

	// Remove from clients list
	ws.Clients.Delete(client.clientName)

	// Update subscriptions, unless close reason is for reconnect testing.
	if ws.Status == 2 {
		ws.muSubscriptions.Lock()
		subscriptions := ws.Subscriptions[client.clientName]
		for i := range subscriptions {
			if subscriptions[i].Status == STATUS_ENABLED {
				subscriptions[i].Status = getStatusFromCloseMessage(closeReason)
				subscriptions[i].ClientConnectedAt = ""
				subscriptions[i].ClientDisconnectedAt = time.Now().UTC().Format(time.RFC3339Nano)
			}
		}
		ws.Subscriptions[client.clientName] = subscriptions
		ws.muSubscriptions.Unlock()
	}

	log.Printf("Disconnected client [%v] with code [%v]", client.clientName, closeReason.code)

	// Print new clients connections list
	ws.printConnections()
}

func (ws *WebSocketServer) printConnections() {
	currentConnections := ""

	for _, client := range ws.Clients.All() {
		currentConnections += client.clientName + ", "
	}

	if currentConnections != "" {
		currentConnections = string(currentConnections[:len(currentConnections)-2])
	}

	log.Printf("[%s] Connections: (%d) [ %s ]", ws.ServerId, ws.Clients.Length(), currentConnections)
}
