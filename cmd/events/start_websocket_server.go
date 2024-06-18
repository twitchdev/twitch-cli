package events

import "github.com/spf13/cobra"

func StartWebsocketServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:        "start-websocket-server",
		Deprecated: `use "twitch event websocket start-server" instead.`,
	}
}
