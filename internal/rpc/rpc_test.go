package rpc_handler

import (
	"encoding/json"
	"net/rpc"
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
)

type rpcTestStruct struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestRPCVerify(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	field1 := "abcd"
	field2 := 1234
	variables := make(map[string]string)
	variables["Test1"] = "123"
	variables["Test2"] = "789"

	r := RPCHandler{
		Port:     44748,
		Handlers: make(map[string]HandlerCallback),
	}

	err := r.StartBackgroundServer()
	a.Nil(err, nil)

	client, err := rpc.DialHTTP("tcp", ":44748")
	a.Nil(err)
	defer client.Close()

	rpcBodyFormatted := rpcTestStruct{
		Field1: field1,
		Field2: field2,
	}
	b, err := json.Marshal(rpcBodyFormatted)
	a.Nil(err)

	args := RPCArgs{
		RPCName:   "Test",
		Body:      string(b),
		Variables: variables,
	}

	var reply RPCArgs
	err = client.Call("RPCHandler.Verify", args, &reply)
	a.Nil(err)

	var formattedReply rpcTestStruct
	err = json.Unmarshal([]byte(reply.Body), &formattedReply)
	a.Nil(err)

	a.Equal(args.RPCName, reply.RPCName)
	a.Equal(rpcBodyFormatted.Field1, formattedReply.Field1)
	a.Equal(rpcBodyFormatted.Field2, formattedReply.Field2)
	a.Equal(args.Variables["Test1"], reply.Variables["Test1"])
	a.Equal(args.Variables["Test2"], reply.Variables["Test2"])

	r.listener.Close()
}
