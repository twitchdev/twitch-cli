package mock_server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"regexp"
	"strconv"
	"strings"
)

var sessionRegex = regexp.MustCompile(`(?P<server_name>.+)_(?P<client_name>.+)`)

type RPCArgs struct {
	Body       string
	ClientName string

	ServerCommand      string
	SubscriptionID     string
	SubscriptionStatus string
	CloseReason        string
}

type CommandResponse struct {
	ResponseCode   int
	AdditionalInfo string
}

const (
	COMMAND_RESPONSE_SUCCESS          int = 0
	COMMAND_RESPONSE_INVALID_CMD      int = 1
	COMMAND_RESPONSE_FAILED_ON_SERVER int = 2
	COMMAND_RESPONSE_MISSING_FLAG     int = 3
)

type WebSocketServerRPC struct{}

// RPC function to execute EventSub events on the WebSocket server.
// May be used to send it to a single client, or all clients.
func (wsrpc *WebSocketServerRPC) RemoteFireEventSub(args *RPCArgs, reply *bool) error {
	server, ok := serverManager.serverList.Get(serverManager.primaryServer)
	if !ok {
		log.Printf("Error on RPC call (RemoteFireEventSub): Primary server not in server list.")
		*reply = false
		return nil
	}

	clientName := args.ClientName
	if sessionRegex.MatchString(clientName) {
		// Users can include the full session_id given in the response. If they do, subtract it to just the client name
		clientName = sessionRegex.FindAllStringSubmatch(clientName, -1)[0][2]
	}

	success := server.HandleRPCEventSubForwarding(args.Body, clientName)

	*reply = success
	return nil
}

// RPC function to attempt to start reconnect testing on the primary WebSocket server
func (wsrpc *WebSocketServerRPC) ServerCommand(args *RPCArgs, reply *CommandResponse) error {
	// Initiate reconnect
	//success := serverManager.HandleRPCServerCommand(args.Body)

	cmd := strings.ToLower(args.ServerCommand)

	if cmd == "reconnect" {
		success := serverManager.HandleRPCCommandReconnect(args.Body)
		if !success {
			*reply = CommandResponse{
				ResponseCode: COMMAND_RESPONSE_FAILED_ON_SERVER,
			}
		}

	} else if cmd == "close" {
		closeCode, err := strconv.Atoi(args.CloseReason)

		if err != nil || args.ClientName == "" || args.CloseReason == "" {
			*reply = CommandResponse{
				ResponseCode: COMMAND_RESPONSE_MISSING_FLAG,
				AdditionalInfo: "Command \"close\" requires flags --client and --reason" +
					"\nThe flag --reason must be one of the number codes defined here:" +
					"\nhttps://dev.twitch.tv/docs/eventsub/websocket-reference/#close-message" +
					"\n\nExample: twitch event websocket close --client=4a1ab390 --reason=4006",
			}

			return nil
		}

		success, failReason := serverManager.HandleRPCCommandClose(args.ClientName, closeCode)
		if success {
			*reply = CommandResponse{
				ResponseCode: COMMAND_RESPONSE_SUCCESS,
			}
		} else {
			*reply = CommandResponse{
				ResponseCode:   COMMAND_RESPONSE_FAILED_ON_SERVER,
				AdditionalInfo: failReason,
			}
		}

	} else if cmd == "subscription" {
		if args.SubscriptionID == "" || !IsValidSubscriptionStatus(args.SubscriptionStatus) {
			*reply = CommandResponse{
				ResponseCode: COMMAND_RESPONSE_MISSING_FLAG,
				AdditionalInfo: "Command \"subscription\" requires flags --status, --subscription, and --client" +
					fmt.Sprintf("\nThe flag --subscription must be the ID of the subscription made at http://localhost:%v/eventsub/subscriptions", serverManager.port) +
					"\nThe flag --status must be one of the non-webhook status options defined here:" +
					"\nhttps://dev.twitch.tv/docs/api/reference/#get-eventsub-subscriptions" +
					"\n\nExample: twitch event websocket subscription --client=4a1ab390 --status=user_removed --subscription=48d3-b9a-f84c",
			}

			return nil
		}

		success, failReason := serverManager.HandleRPCCommandSubscription(args.SubscriptionStatus, args.SubscriptionID)
		if success {
			*reply = CommandResponse{
				ResponseCode: COMMAND_RESPONSE_SUCCESS,
			}
		} else {
			*reply = CommandResponse{
				ResponseCode:   COMMAND_RESPONSE_FAILED_ON_SERVER,
				AdditionalInfo: failReason,
			}
		}

	} else {
		*reply = CommandResponse{
			ResponseCode: COMMAND_RESPONSE_INVALID_CMD,
		}
	}

	return nil
}

func StartRPCListener() error {
	rpc.Register(new(WebSocketServerRPC))
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":44747")
	if err != nil {
		return err
	}
	go http.Serve(l, nil)

	return nil
}
