package mock_ws

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"regexp"
)

var sessionRegex = regexp.MustCompile(`(?P<server_name>.+)_(?P<client_name>.+)`)

type RPCArgs struct {
	Body     string
	ClientID string
}

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

	clientId := args.ClientID
	if sessionRegex.MatchString(clientId) {
		// Users can include the full session_id given in the response. If they do, subtract it to just the client name
		clientId = sessionRegex.FindAllStringSubmatch(clientId, -1)[0][2]
	}

	success := server.HandleRPCEventSubForwarding(args.Body, clientId)

	*reply = success
	return nil
}

// RPC function to attempt to start reconnect testing on the primary WebSocket server
func (wsrpc *WebSocketServerRPC) ServerCommand(args *RPCArgs, reply *bool) error {
	// Initiate reconnect
	success := serverManager.HandleRPCServerCommand(args.Body)

	*reply = success
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
