package websocket

import (
	"fmt"
	"net/rpc"

	"github.com/fatih/color"
	"github.com/twitchdev/twitch-cli/internal/events/websocket/mock_server"
	rpc_handler "github.com/twitchdev/twitch-cli/internal/rpc"
)

type WebsocketCommandParameters struct {
	Client             string
	Subscription       string
	SubscriptionStatus string
	CloseReason        string
}

func ForwardWebsocketCommand(cmd string, p WebsocketCommandParameters) error {
	client, err := rpc.DialHTTP("tcp", ":44747")
	if err != nil {
		return fmt.Errorf("Failed to dial RPC handler for WebSocket server. Is it online?\nError: %v", err.Error())
	}

	var reply rpc_handler.RPCResponse

	rpcName := mock_server.ResolveRPCName(cmd)
	if rpcName == "" {
		return fmt.Errorf("Invalid websocket command")
	}

	// Command line flags to be passed with the command
	// Add them all, as it wont hurt anything if they're not relevant
	variables := make(map[string]string)
	variables["ClientName"] = p.Client
	variables["SubscriptionID"] = p.Subscription
	variables["SubscriptionStatus"] = p.SubscriptionStatus
	variables["CloseReason"] = p.CloseReason

	args := &rpc_handler.RPCArgs{
		RPCName:   rpcName,
		Variables: variables,
	}

	err = client.Call("RPCHandler.ExecuteGenericRPC", args, &reply)

	if err != nil {
		return fmt.Errorf("Failed to call RPC method RPCHandler.ExecuteGenericRPC: %v", err.Error())
	}

	switch reply.ResponseCode {
	case mock_server.COMMAND_RESPONSE_SUCCESS:
		color.New().Add(color.FgGreen).Println(fmt.Sprintf("✔ Forwarded for use in mock EventSub WebSocket server"))
		return nil

	case mock_server.COMMAND_RESPONSE_FAILED_ON_SERVER:
		return fmt.Errorf(
			color.New().Add(color.FgRed).Sprintln(fmt.Sprintf("✗ EventSub WebSocket server failed to process command:\n%v", reply.DetailedInfo)),
		)

	case mock_server.COMMAND_RESPONSE_MISSING_FLAG:
		return fmt.Errorf(
			color.New().Add(color.FgRed).Sprintln(fmt.Sprintf("✗ Command rejected for invalid flags:\n%v", reply.DetailedInfo)),
		)

	case mock_server.COMMAND_RESPONSE_INVALID_CMD:
		return fmt.Errorf("Invalid websocket sub-command: %v", cmd)

	}

	return fmt.Errorf("RPCHandler experienced unexpected response code: %v", reply.ResponseCode)
}
