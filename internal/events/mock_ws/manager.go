package mock_ws

import (
	"os"
	"os/signal"

	"github.com/twitchdev/twitch-cli/internal/events/mock_ws/mock_ws_server"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type ServerManager struct {
	serverList    *mock_ws_server.Servers
	primaryServer string
}

var serverManager ServerManager

func StartWebsocketServers(enableDebug bool, port int) {
	// Initialize server manager
	serverManager = ServerManager{
		serverList: &mock_ws_server.Servers{
			Servers: make(map[string]*mock_ws_server.WebSocketServer),
		},
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	// Start initial server
	initialServer := &mock_ws_server.WebSocketServer{
		ServerId: util.RandomGUID()[:8],
		Status:   2,
		Clients: &mock_ws_server.Clients{
			Clients: make(map[string]*mock_ws_server.Client),
		},
	}
	initialServer.StartServer(port, enableDebug)
	serverManager.serverList.Put(initialServer.ServerId, initialServer)
	serverManager.primaryServer = initialServer.ServerId

	// Initalize RPC handler, to accept EventSub transports
	StartRPCListener()

	<-stop // Wait for Ctrl + C
}
