package mock_ws

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type RPCArgs struct {
	Body     string
	ClientID string
}

type WebSocketServerRPC struct{}

func (wsrpc *WebSocketServerRPC) RemoteFireEventSub(args *RPCArgs, reply *bool) error {
	server, ok := serverManager.serverList.Get(serverManager.primaryServer)
	if !ok {
		log.Printf("Error on RPC call: Primary server not in server list.")
		*reply = false
		return nil
	}
	success := server.HandleRPCForwarding(args.Body, args.ClientID)

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
