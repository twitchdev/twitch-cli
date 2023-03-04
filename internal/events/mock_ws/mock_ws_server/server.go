package mock_ws_server

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
	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/generate"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

// Minimum time between messages before the server disconnects a client.
const KEEPALIVE_TIMEOUT_SECONDS = 10

var upgrader = websocket.Upgrader{}
var debugEnabled = false

type WebSocketServer struct {
	ServerId string   // Int representing the ID of the server
	Status   int      // 0 = shut down; 1 = shutting down; 2 = online
	Clients  *Clients // All connected clients
}

func (ws *WebSocketServer) StartServer(port int, debug bool) {
	debugEnabled = debug

	log.Printf("Attempting to start WebSocket server [%v] on port %v", ws.ServerId, port)

	m := http.NewServeMux()

	// Connect to SQLite database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err.Error())
		return
	}

	// Generate database if this is the first run.
	firstTime := db.IsFirstRun()
	if firstTime {
		err := generate.Generate(25)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Register URL handler
	m.HandleFunc("/ws", ws.wsPageHandler)

	// Allow exit with Ctrl + C
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

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
}

func (ws *WebSocketServer) wsPageHandler(w http.ResponseWriter, r *http.Request) {
	// This next line is required to disable CORS checking. No sense in caring in a test environment.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
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
		clientId: util.RandomGUID()[:8],
		conn:     conn,
	}

	// Disconnect the user if the server is in reconnect phase
	if ws.Status == 1 || ws.Status == 0 {
		log.Printf("New client trying to connect while websocket server in reconnect phase. Disconnecting them.")
		// TODO: Change this based on what was given in #eventsub
		client.CloseWithReason(closeNetworkError)
		// No handleClientConnectionClose because client is not in clients list, and chan loop was not set up yet.
	}

	// Add to the client connections list
	ws.Clients.Put(client.clientId, client)

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
					ID:                      ws.ServerId,
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
	mustSubscribeTicker := time.NewTimer(10 * time.Second)
	go func() {
		select {
		case <-mustSubscribeTicker.C:
			mustSubscribeTicker.Stop()
			log.Printf("No subscriptions whomp whomp")
			// TODO: Check for subscriptions. If none, disconnect.
		}
	}()

	// Set up ping/pong and keepalive handling
	ticker := time.NewTicker(5 * time.Second)
	client.pingKeepAliveLoopChan = make(chan struct{})
	go func() {
		// Set pong handler. Resets the read deadline when pong is received.
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(time.Second * KEEPALIVE_TIMEOUT_SECONDS))
			return nil
		})

		keepAliveNext := false

		for {
			select {
			case <-client.pingKeepAliveLoopChan:
				ticker.Stop()
				return
			case <-ticker.C:
				// Send keepalive every 10 seconds (two 5 second loops)
				if !keepAliveNext {
					keepAliveNext = true
				} else {
					keepAliveNext = false
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

					if debugEnabled {
						log.Printf("Sent session_keepalive to client [%s]", client.clientId)
					}
				}

				// Send ping
				err := client.SendMessage(websocket.PingMessage, []byte{})
				if err != nil {
					client.CloseWithReason(closeClientFailedPingPong)
					ws.handleClientConnectionClose(client)
				}
			}
		}
	}()

	// Wait for message
	for {
		// Reset timeout upon every message, no matter what it is.
		client.conn.SetReadDeadline(time.Now().Add(time.Second * KEEPALIVE_TIMEOUT_SECONDS))

		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read err [%v]: %v", client.clientId, err)
			client.CloseWithReason(closeNetworkError)
			ws.handleClientConnectionClose(client)
			break
		}

		log.Printf("Disconnecting client [%v] due to received inbound traffic.\nMessage[%d]: %s", client.clientId, mt, message)
		client.CloseWithReason(closeClientSentInboundTraffic)
		ws.handleClientConnectionClose(client)
		break
	}
}

func (ws *WebSocketServer) HandleRPCForwarding(eventsubBody string, clientId string) bool {
	// If --client is used, make sure the client exists
	if clientId != "" {
		_, ok := ws.Clients.Get(strings.ToLower(clientId))
		if !ok {
			log.Printf("Error executing remote triggered EventSub: Client [%v] does not exist", clientId)
			return false
		}
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

		// Change transport around
		eventObj.Subscription.Transport.Callback = ""
		eventObj.Subscription.Transport.SessionID = client.clientId

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

func (ws *WebSocketServer) handleClientConnectionClose(client *Client) {
	// Prevent further looping
	close(client.pingKeepAliveLoopChan)

	// Remove from clients list
	ws.Clients.Delete(client.clientId)

	log.Printf("Disconnected client [%v]", client.clientId)

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
