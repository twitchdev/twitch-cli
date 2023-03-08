package mock_ws_server

import (
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
	ServerId     string // Int representing the ID of the server
	ReconnectUrl string // Server's url for people to connect to. Used for messaging in reconnect testing
	DebugEnabled bool   // Display debug messages; --debug
	StrictMode   bool   // Force stricter production-like qualities; --strict
	Upgrader     websocket.Upgrader

	Clients   *Clients   // All connected clients
	muClients sync.Mutex // Muted for WebSocketServer.Clients

	Status   int        // 0 = shut down; 1 = shutting down; 2 = online
	muStatus sync.Mutex // Mutex for WebSocketServer.Status

	ReconnectClients   map[string]*[]string // Clients that were part of the last server
	muReconnectClients sync.Mutex
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
		clientId:                 util.RandomGUID()[:8],
		conn:                     conn,
		clientConnectedTimestamp: connectedAtTimestamp,

		keepAliveChanOpen: false,
		pingChanOpen:      false,
	}

	// Disconnect the user if the server is in reconnect phase
	ws.muStatus.Lock()
	if ws.Status != 2 {
		// This is the closest we can get to the production environment, as there's no way to route people to a shutting down server
		log.Printf("New client trying to connect while websocket server in reconnect phase. Disconnecting them.")
		client.CloseDirty()
		// No handleClientConnectionClose because client is not in clients list, and chan loop was not set up yet.
	}

	// TODO: Check if user is connected to the old server, and if they are then disconnect them from old server with close frame 4004

	ws.muClients.Lock()
	// Add to the client connections list
	ws.Clients.Put(client.clientId, client)
	ws.muClients.Unlock()

	// This is put after ws.Clients.Put to make sure the client gets included in the list before InitiateRestart() kicks everyone out
	// Avoids any possible rare edge cases. This ain't production but I can still be safe :)
	ws.muStatus.Unlock()

	log.Printf("Client connected [%v]", client.clientId)
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
					ID:                      fmt.Sprintf("%v_%v", ws.ServerId, client.clientId),
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
	client.mustSubscribeTimer = time.NewTimer(10 * time.Second)
	go func() {
		select {
		case <-client.mustSubscribeTimer.C:
			log.Printf("No subscriptions whomp whomp -- TODO")
			// TODO: Check for subscriptions. If none, disconnect.
		}
	}()

	// Set up ping/pong and keepalive handling
	client.keepAliveTimer = time.NewTicker(10 * time.Second)
	client.pingTimer = time.NewTicker(5 * time.Second)
	client.keepAliveLoopChan = make(chan struct{})
	client.pingLoopChan = make(chan struct{})
	client.keepAliveChanOpen = true
	client.pingChanOpen = true
	go func() {
		// Set pong handler. Resets the read deadline when pong is received.
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(time.Second * KEEPALIVE_TIMEOUT_SECONDS))
			return nil
		})

		for {
			select {
			case <-client.keepAliveLoopChan:
				client.keepAliveTimer.Stop()

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
					log.Printf("Sent session_keepalive to client [%s]", client.clientId)
				}

			case <-client.pingLoopChan:
				client.pingTimer.Stop()

			case <-client.pingTimer.C: // Send ping
				err := client.SendMessage(websocket.PingMessage, []byte{})
				if err != nil {
					ws.muClients.Lock()
					client.CloseWithReason(closeClientFailedPingPong)
					ws.handleClientConnectionClose(client, closeClientFailedPingPong)
					ws.muClients.Unlock()
				}

				if ws.DebugEnabled {
					log.Printf("Sent pong to client [%s]", client.clientId)
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
			log.Printf("read err [%v]: %v", client.clientId, err)

			ws.muClients.Lock()
			client.CloseWithReason(closeNetworkError)
			ws.handleClientConnectionClose(client, closeNetworkError)
			ws.muClients.Unlock()
			break
		}

		if ws.Status == 2 { // Only care about this when the server is running
			log.Printf("Disconnecting client [%v] due to received inbound traffic.\nMessage[%d]: %s", client.clientId, mt, message)

			ws.muClients.Lock()
			client.CloseWithReason(closeClientSentInboundTraffic)
			ws.handleClientConnectionClose(client, closeClientSentInboundTraffic)
			ws.muClients.Unlock()
		}

		break
	}
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
		reconnectMsg, _ := json.Marshal(
			ReconnectMessage{
				Metadata: MessageMetadata{
					MessageID:        util.RandomGUID(),
					MessageType:      "session_reconnect",
					MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
				},
				Payload: ReconnectMessagePayload{
					Session: ReconnectMessagePayloadSession{
						ID:                      fmt.Sprintf("%v_%v", ws.ServerId, client.clientId),
						Status:                  "reconnecting",
						KeepaliveTimeoutSeconds: nil,
						ReconnectUrl:            ws.ReconnectUrl,
						ConnectedAt:             client.clientConnectedTimestamp,
					},
				},
			},
		)

		err := client.SendMessage(websocket.TextMessage, reconnectMsg)
		if err != nil {
			log.Printf("Error building session_reconnect JSON for client [%v]: %v", client.clientId, err.Error())
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

func (ws *WebSocketServer) HandleRPCEventSubForwarding(eventsubBody string, clientId string) bool {
	// If --client is used, make sure the client exists
	if clientId != "" {
		_, ok := ws.Clients.Get(strings.ToLower(clientId))
		if !ok {
			log.Printf("Error executing remote triggered EventSub: Client [%v] does not exist on server [%v]", clientId, ws.ServerId)
			return false
		}
	}

	if ws.Clients.Length() == 0 {
		log.Printf("Warning for remote triggered EventSub: No clients in server [%v]", ws.ServerId)
		return false
	}

	for _, client := range ws.Clients.All() {
		if clientId != "" && !strings.EqualFold(strings.ToLower(clientId), clientId) {
			// When --client is used, only send to that client
			continue
		}

		// Convert to struct for editing
		eventObj := models.EventsubResponse{}
		err := json.Unmarshal([]byte(eventsubBody), &eventObj)
		if err != nil {
			log.Printf("Error reading JSON forwarded from EventSub. Currently on client [%v]: %v\nRaw: %v", client.clientId, err.Error(), eventsubBody)
			return false
		}

		// Change payload's subscription.transport.session_id to contain the correct Session ID
		eventObj.Subscription.Transport.SessionID = fmt.Sprintf("%v_%v", ws.ServerId, client.clientId)

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
			log.Printf("Error building JSON for client [%v]: %v", client.clientId, err.Error())
			return false
		}

		client.SendMessage(websocket.TextMessage, notificationMsg)
		log.Printf("Sent [%v / %v] to client [%v]", eventObj.Subscription.Type, eventObj.Subscription.Version, client.clientId)
	}

	return true
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
	ws.Clients.Delete(client.clientId)

	log.Printf("Disconnected client [%v] with code [%v]", client.clientId, closeReason.code)

	// Print new clients connections list
	ws.printConnections()
}

func (ws *WebSocketServer) printConnections() {
	currentConnections := ""

	for _, client := range ws.Clients.All() {
		currentConnections += client.clientId + ", "
	}

	if currentConnections != "" {
		currentConnections = string(currentConnections[:len(currentConnections)-2])
	}

	log.Printf("[%s] Connections: (%d) [ %s ]", ws.ServerId, ws.Clients.Length(), currentConnections)
}
