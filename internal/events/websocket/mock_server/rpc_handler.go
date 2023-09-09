package mock_server

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/gorilla/websocket"
	rpc "github.com/twitchdev/twitch-cli/internal/rpc"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var sessionRegex = regexp.MustCompile(`(?P<server_name>.+)_(?P<client_name>.+)`)

const (
	COMMAND_RESPONSE_SUCCESS          int = 0
	COMMAND_RESPONSE_INVALID_CMD      int = 1
	COMMAND_RESPONSE_FAILED_ON_SERVER int = 2
	COMMAND_RESPONSE_MISSING_FLAG     int = 3
)

// Resolves console commands to their RPC names defined in the server manager
func ResolveRPCName(cmd string) string {
	if cmd == "reconnect" {
		return "EventSubWebSocketReconnect"
	} else if cmd == "close" {
		return "EventSubWebSocketCloseClient"
	} else if cmd == "subscription" {
		return "EventSubWebSocketSubscription"
	} else {
		return ""
	}
}

// $ twitch event websocket reconnect
func RPCReconnectHandler(args rpc.RPCArgs) rpc.RPCResponse {
	// Initiate reconnect testing
	log.Printf("Initiating reconnect testing...")

	if serverManager.reconnectTesting {
		log.Printf("Error on RPC call (EventSubWebSocketReconnect): Cannot execute reconnect testing while its already in progress. Discarding duplicate reconnect command.")
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: "Error on RPC call (EventSubWebSocketReconnect): Cannot execute reconnect testing while its already in progress. Discarding duplicate reconnect command.",
		}
	}

	// Find current primary server
	originalPrimaryServer, ok := serverManager.serverList.Get(serverManager.primaryServer)
	if !ok {
		log.Printf("Error on RPC call (EventSubWebSocketReconnect): Primary server not in server list.")
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: "Error on RPC call (EventSubWebSocketReconnect): Primary server not in server list.",
		}
	}

	serverManager.reconnectTesting = true

	// Get the list of reconnect clients ready
	reconnectClients := originalPrimaryServer.GetCurrentSubscriptionsForReconnect()

	// Spin up new server
	newServer := &WebSocketServer{
		ServerId: util.RandomGUID()[:8],
		Status:   2,
		Clients: &util.List[Client]{
			Elements: make(map[string]*Client),
		},
		Upgrader:         websocket.Upgrader{},
		DebugEnabled:     serverManager.debugEnabled,
		Subscriptions:    make(map[string][]Subscription),
		StrictMode:       serverManager.strictMode,
		ReconnectClients: reconnectClients,
	}
	serverManager.serverList.Put(newServer.ServerId, newServer)

	// Switch manager's primary server to new one
	// Doing this before sending the reconnect messages emulates the Twitch's production load balancer, which will never send to servers shutting down.
	serverManager.primaryServer = newServer.ServerId

	// Notify primary server to restart (includes not accepting new clients)
	// This is in a goroutine so it doesn't hang the reconnect command
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

		log.Printf("Reconnect testing successful. Primary server is now [%v]\nYou may now execute reconnect testing again.", serverManager.primaryServer)
	}()

	return rpc.RPCResponse{
		ResponseCode: COMMAND_RESPONSE_SUCCESS,
	}
}

// $ twitch event trigger <event> --transport=websocket
func RPCFireEventSubHandler(args rpc.RPCArgs) rpc.RPCResponse {
	server, ok := serverManager.serverList.Get(serverManager.primaryServer)
	if !ok {
		log.Printf("Error on RPC call (EventSubWebSocketForwardEvent): Primary server not in server list.")
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: "Primary server not in server list.",
		}
	}

	clientName := args.Variables["ClientName"]
	if sessionRegex.MatchString(clientName) {
		// Users can include the full session_id given in the response. If they do, subtract it to just the client name
		clientName = sessionRegex.FindAllStringSubmatch(clientName, -1)[0][2]
	}

	success, failMsg := server.HandleRPCEventSubForwarding(args.Body, clientName)

	if success {
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_SUCCESS,
		}
	} else {
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: failMsg,
		}
	}
}

// $ twitch event websocket close
func RPCCloseHandler(args rpc.RPCArgs) rpc.RPCResponse {
	closeCode, err := strconv.Atoi(args.Variables["CloseReason"])

	if err != nil || args.Variables["ClientName"] == "" || args.Variables["CloseReason"] == "" {
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_MISSING_FLAG,
			DetailedInfo: "Command \"close\" requires flags --session and --reason" +
				"\nThe flag --reason must be one of the number codes defined here:" +
				"\nhttps://dev.twitch.tv/docs/eventsub/websocket-reference/#close-message" +
				"\n\nExample: twitch event websocket close --session=e411cc1e_a2613d4e --reason=4006",
		}
	}

	clientName := args.Variables["ClientName"]

	if serverManager.reconnectTesting {
		log.Printf("Error on RPC call (EventSubWebSocketCloseClient): Could not activate while reconnect testing is active.")
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: "Cannot activate this command while reconnect testing is active.",
		}
	}

	server, ok := serverManager.serverList.Get(serverManager.primaryServer)
	if !ok {
		log.Printf("Error on RPC call (EventSubWebSocketCloseClient): Primary server not in server list.")
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: "Primary server not in server list.",
		}
	}

	cn := clientName
	if sessionRegex.MatchString(clientName) {
		// Client name given was formatted as <server_id>_<client_name>. We must extract it
		sessionRegexExec := sessionRegex.FindAllStringSubmatch(clientName, -1)
		cn = sessionRegexExec[0][2]
	}

	server.muClients.Lock()

	client, ok := server.Clients.Get(cn)
	if !ok {
		server.muClients.Unlock()
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: "Client [" + cn + "] does not exist on WebSocket server.",
		}
	}

	closeMessage := GetCloseMessageFromCode(closeCode)
	if closeMessage == nil {
		server.muClients.Unlock()
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: fmt.Sprintf("Close code [%v] not supported.", closeCode),
		}
	}

	server.muClients.Unlock()

	client.CloseWithReason(closeMessage)
	server.handleClientConnectionClose(client, closeMessage)

	log.Printf("RPC instructed to close client [%v] with code [%v]", clientName, closeCode)

	return rpc.RPCResponse{
		ResponseCode: COMMAND_RESPONSE_SUCCESS,
	}
}

// $ twitch event websocket subscription
func RPCSubscriptionHandler(args rpc.RPCArgs) rpc.RPCResponse {
	if args.Variables["SubscriptionID"] == "" || !IsValidSubscriptionStatus(args.Variables["SubscriptionStatus"]) {
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_MISSING_FLAG,
			DetailedInfo: "Command \"subscription\" requires flags --status, --subscription, and --session" +
				fmt.Sprintf("\nThe flag --subscription must be the ID of the subscription made at %v://%v:%v/eventsub/subscriptions", serverManager.protocolHttp, serverManager.ip, serverManager.port) +
				"\nThe flag --status must be one of the non-webhook status options defined here:" +
				"\nhttps://dev.twitch.tv/docs/api/reference/#get-eventsub-subscriptions" +
				"\n\nExample: twitch event websocket subscription --status=user_removed --subscription=82a855-fae8-93bff0",
		}
	}

	if serverManager.reconnectTesting {
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: "Cannot activate this command while reconnect testing is active.",
		}
	}

	server, ok := serverManager.serverList.Get(serverManager.primaryServer)
	if !ok {
		log.Printf("Error on RPC call (EventSubWebSocketSubscription): Primary server not in server list.")
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: "Primary server not in server list.",
		}
	}

	server.muSubscriptions.Lock()
	found := false
	for client, clientSubscriptions := range server.Subscriptions {
		if found {
			break
		}

		for i, sub := range clientSubscriptions {
			if sub.SubscriptionID == args.Variables["SubscriptionID"] {
				found = true

				server.Subscriptions[client][i].Status = args.Variables["SubscriptionStatus"]
				if args.Variables["SubscriptionStatus"] == STATUS_ENABLED {
					server.Subscriptions[client][i].DisabledAt = nil
				} else {
					tNow := util.GetTimestamp()
					server.Subscriptions[client][i].DisabledAt = &tNow
				}
				break
			}
		}
	}
	server.muSubscriptions.Unlock()

	if !found {
		return rpc.RPCResponse{
			ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			DetailedInfo: fmt.Sprintf("Subscription ID [%v] does not exist", args.Variables["SubscriptionID"]),
		}
	}

	return rpc.RPCResponse{
		ResponseCode: COMMAND_RESPONSE_SUCCESS,
	}
}
