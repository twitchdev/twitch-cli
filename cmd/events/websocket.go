package events

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/twitchdev/twitch-cli/internal/events/websocket"
	"github.com/twitchdev/twitch-cli/internal/events/websocket/mock_server"
)

var (
	wsDebug          bool
	wsStrict         bool
	wsClient         string
	wsSubscription   string
	wsStatus         string
	wsReason         string
	wsServerIP       string
	wsServerPort     int
	wsSSL            bool
	wsFeatureEnabled bool
)

func WebsocketCommand() (command *cobra.Command) {
	command = &cobra.Command{
		Use:   "websocket [action]",
		Short: `Executes actions regarding the mock EventSub WebSocket server. See "twitch event websocket --help" for usage info.`,
		Long:  "Executes actions regarding the mock EventSub WebSocket server.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  websocketCmdRun,
		Example: `  twitch event websocket start-server
	  twitch event websocket reconnect
	  twitch event websocket close --session=e411cc1e_a2613d4e --reason=4006
	  twitch event websocket subscription --status=user_removed --subscription=82a855-fae8-93bff0
	  twitch event websocket keepalive --session=e411cc1e_a2613d4e --enabled=false`,
		Aliases: []string{
			"websockets",
			"ws",
			"wss",
		},
	}

	// flags for start-server
	command.Flags().StringVar(&wsServerIP, "ip", "127.0.0.1", "Defines the ip that the mock EventSub websocket server will bind to.")
	command.Flags().IntVarP(&wsServerPort, "port", "p", 8080, "Defines the port that the mock EventSub websocket server will run on.")
	command.Flags().BoolVar(&wsSSL, "ssl", false, "Enables SSL for EventSub websocket server (wss) and EventSub mock subscription server (https).")
	command.Flags().BoolVar(&wsDebug, "debug", false, "Set on/off for debug messages for the EventSub WebSocket server.")
	command.Flags().BoolVarP(&wsStrict, "require-subscription", "S", false, "Requires subscriptions for all events, and activates 10 second subscription requirement.")

	// flags for everything else
	command.Flags().StringVarP(&wsClient, "session", "s", "", "WebSocket client/session to target with your server command. Used in multiple commands.")
	command.Flags().StringVar(&wsSubscription, "subscription", "", `Subscription to target with your server command. Used with "websocket subscription".`)
	command.Flags().StringVar(&wsStatus, "status", "", `Changes the status of an existing subscription. Used with "websocket subscription".`)
	command.Flags().StringVar(&wsReason, "reason", "", `Sets the close reason when sending a Close message to the client. Used with "websocket close".`)
	command.Flags().BoolVar(&wsFeatureEnabled, "enabled", false, "Sets on/off for the specified feature.")

	return
}

func websocketCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		return fmt.Errorf("")
	}

	if args[0] == "start-server" || args[0] == "start" {
		log.Printf("Attempting to start WebSocket server on %v:%v", wsServerIP, wsServerPort)
		log.Printf("`Ctrl + C` to exit mock WebSocket servers.")
		mock_server.StartWebsocketServer(wsDebug, wsServerIP, wsServerPort, wsSSL, wsStrict)
	} else {
		// Forward all other commands via RPC
		err := websocket.ForwardWebsocketCommand(args[0], websocket.WebsocketCommandParameters{
			Client:             wsClient,
			Subscription:       wsSubscription,
			SubscriptionStatus: wsStatus,
			CloseReason:        wsReason,
			FeatureEnabled:     wsFeatureEnabled,
		})

		return err
	}

	return nil
}
