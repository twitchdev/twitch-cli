package websocket

import (
	"fmt"
	"net/rpc"

	"github.com/fatih/color"
	"github.com/twitchdev/twitch-cli/internal/events/websocket/mock_server"
)

type WebsocketCommandParameters struct {
	Client             string
	Subscription       string
	SubscriptionStatus string
	CloseReason        string
}

func ForwardWebsocketCommand(cmd string, p WebsocketCommandParameters) {
	client, err := rpc.DialHTTP("tcp", ":44747")
	if err != nil {
		println("Failed to dial RPC handler for WebSocket server. Is it online?")
		println("Error: " + err.Error())
		return
	}

	var reply mock_server.CommandResponse

	args := &mock_server.RPCArgs{
		ServerCommand:      cmd,
		ClientName:         p.Client,
		SubscriptionID:     p.Subscription,
		SubscriptionStatus: p.SubscriptionStatus,
		CloseReason:        p.CloseReason,
	}

	err = client.Call("WebSocketServerRPC.ServerCommand", args, &reply)

	switch reply.ResponseCode {
	case mock_server.COMMAND_RESPONSE_SUCCESS:
		color.New().Add(color.FgGreen).Println(fmt.Sprintf("✔ Forwarded for use in mock EventSub WebSocket server"))

	case mock_server.COMMAND_RESPONSE_FAILED_ON_SERVER:
		color.New().Add(color.FgRed).Println(fmt.Sprintf("✗ EventSub WebSocket server failed to process command:\n%v", reply.AdditionalInfo))

	case mock_server.COMMAND_RESPONSE_MISSING_FLAG:
		color.New().Add(color.FgRed).Println(fmt.Sprintf("✗ Command rejected for invalid flags:\n%v", reply.AdditionalInfo))

	case mock_server.COMMAND_RESPONSE_INVALID_CMD:
		println("Invalid websocket sub-command: " + cmd)

	}
}
