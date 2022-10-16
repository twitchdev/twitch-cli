package mock_wss_server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/websocket"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/generate"
	"github.com/twitchdev/twitch-cli/internal/util"
)

/* Minimum time between messages before the server disconnects a client. AKA, "timeout period"
 * This is a const for now, but it may need to be a flag in the future. */
const MINIMUM_MESSAGE_FREQUENCY_SECONDS = 10

//var connectionURL string

var upgrader = websocket.Upgrader{}
var debug = false

// List of websocket servers. Limited to 2 for now, as allowing more would require rewriting reconnect stuff.
var wsServers = [2]*WebsocketServer{}

type WebsocketServer struct {
	serverId          int                   // Int representing the ID of the server (0, 1, 2, ...)
	websocketId       string                // UUID of the websocket. Used for subscribing via EventSub
	connectionUrl     string                // URL used to connect to the websocket. Used for reconnect messages
	connections       []WebsocketConnection // Current clients connected to this websocket
	deactivatedStatus bool                  // Boolean used for preventing connections/messages during deactivation; Used for reconnect testing
}

type WebsocketConnection struct {
	clientId             string
	Conn                 *websocket.Conn
	connectedAtTimestamp string
}

func eventsubHandle(w http.ResponseWriter, r *http.Request) {
	// This next line is required to disable CORS checking.
	// If we weren't returning "true", the the /debug page on the default port wouldn't be able to make calls to port+1
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("[[websocket upgrade err]] ", err)
		return
	}
	defer conn.Close()

	// Connection sucessful. WebSocket is open.

	serverId, _ := strconv.Atoi(r.Context().Value("serverId").(string))
	wsSrv := wsServers[serverId]

	// Or is it? Check for websocket set to deactivated (due to reconnect), and kick them out if so
	if wsSrv.deactivatedStatus {
		log.Printf("Client trying to connect while websocket in reconnect timeout phase. Disconnecting them.")
		conn.Close()
		return
	}

	// TODO: Decline websocket if it reached 100 connections from the same application access token

	// RFC3339Nano = "2022-10-04T12:38:15.548912638Z07" ; This is used by Twitch in production
	connectedAtTimestamp := time.Now().UTC().Format(time.RFC3339Nano)
	conn.SetReadDeadline(time.Now().Add(time.Second * MINIMUM_MESSAGE_FREQUENCY_SECONDS))

	// Add to websocket connection list.
	wsSrv.connections = append(
		wsSrv.connections,
		WebsocketConnection{
			clientId:             util.RandomGUID(),
			Conn:                 conn,
			connectedAtTimestamp: connectedAtTimestamp,
		},
	)
	printConnections(*wsSrv)

	// Send "websocket_welcome" message
	welcomeMsg, _ := json.Marshal(
		WelcomeMessage{
			Metadata: MessageMetadata{
				MessageID:        util.RandomGUID(),
				MessageType:      "websocket_welcome",
				MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
			},
			Payload: WelcomeMessagePayload{
				Websocket: WelcomeMessagePayloadWebsocket{
					ID:                             wsSrv.websocketId,
					Status:                         "connected",
					MinimumMessageFrequencySeconds: MINIMUM_MESSAGE_FREQUENCY_SECONDS,
					ConnectedAt:                    connectedAtTimestamp,
				},
			},
		},
	)
	conn.WriteMessage(1, welcomeMsg)

	if debug {
		log.Printf("[DEBUG] Write: %s", welcomeMsg)
	}

	// TODO: Read messages
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * MINIMUM_MESSAGE_FREQUENCY_SECONDS))
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			onCloseConnection(*wsSrv, conn)
			break
		}
		log.Printf("recv: [%d] %s", mt, message)
		/*err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			onCloseConnection(websocketId)
			break
		}*/
	}
}

func onCloseConnection(wsSrv WebsocketServer, conn *websocket.Conn) {
	// Remove from list
	c := 0
	for i := 0; i < len(wsSrv.connections); i++ {
		if wsSrv.connections[i].Conn == wsSrv.connections[i].Conn {
			log.Printf("Disconnected websocket %s", wsSrv.connections[i].clientId)
			c = i
			break
		}
	}
	wsSrv.connections = append(wsSrv.connections[:c], wsSrv.connections[c+1:]...)

	printConnections(wsSrv)
}

func printConnections(wsSrv WebsocketServer) {
	currentConnections := ""
	for _, s := range wsSrv.connections {
		currentConnections += s.clientId + ", "
	}
	if currentConnections != "" {
		currentConnections = string(currentConnections[:len(currentConnections)-2])
	}
	log.Printf("[Server %v] Connections: (%d) [ %s ]", wsSrv.serverId, len(wsSrv.connections), currentConnections)
}

func activateReconnectTest(server http.Server, ctx context.Context) {
	timer := 30 // 30 seconds, as used by Twitch

	serverId, _ := strconv.Atoi(ctx.Value("serverId").(string))
	wsSrv := wsServers[serverId]

	log.Printf("Terminating server %v with reconnect notice. Disallowing connections in 30 seconds.", serverId)

	var wsAltSrv *WebsocketServer
	if serverId == 0 {
		wsAltSrv = wsServers[1]
	} else {
		wsAltSrv = wsServers[0]
	}

	// Stop processing new messages
	wsSrv.deactivatedStatus = true     // This server
	wsAltSrv.deactivatedStatus = false // Other server; We gotta turn it on to accept connections and whatnot

	log.Printf("Connections: %v", len(wsSrv.connections))

	// Send reconnect notices
	for _, c := range wsSrv.connections {
		reconnectMsg, _ := json.Marshal(
			ReconnectMessage{
				Metadata: MessageMetadata{
					MessageID:        util.RandomGUID(),
					MessageType:      "websocket_reconnect",
					MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
				},
				Payload: ReconnectMessagePayload{
					Websocket: ReconnectMessagePayloadWebsocket{
						ID:                             wsSrv.websocketId,
						Status:                         "reconnecting",
						MinimumMessageFrequencySeconds: MINIMUM_MESSAGE_FREQUENCY_SECONDS,
						Url:                            wsSrv.connectionUrl,
						ConnectedAt:                    c.connectedAtTimestamp,
						ReconnectingAt:                 time.Now().UTC().Add(time.Second * 30).Format(time.RFC3339Nano),
					},
				},
			},
		)

		err := c.Conn.WriteMessage(1, reconnectMsg) // Ignore err here because we're shutting down so it doesn't really matter
		if err != nil {
			log.Printf("ERROR (clientId %v): %v", c.clientId, err)
		} else {
			log.Printf("Sent reconnect notice to %v", c.clientId)
		}
	}

	log.Printf("Reconnect notices sent for server %v", serverId)
	log.Printf("Server \"not accepting connections\" status: [Server 0: %v, Server 1: %v]", wsServers[0].deactivatedStatus, wsServers[1].deactivatedStatus)
	log.Printf("Use this URL for connections: %v", wsAltSrv.connectionUrl)

	// TODO: Transfer subscriptions to the other websocket server.

	// Wait 30 seconds to close out, just like Twitch production EventSub websockets
	select {
	case <-time.After(time.Second * time.Duration(timer)):
		for _, c := range wsSrv.connections {
			c.Conn.Close()
		}

		if debug {
			log.Printf("[DEBUG] Resetting websocket ID on server %v", wsSrv.websocketId)
		}
		wsSrv.websocketId = util.RandomGUID() // Change websocket ID
	}
}

func StartServer(port int, enableDebug bool, reconnectTestTimer int) {
	debug = enableDebug

	m := http.NewServeMux()

	// Server IDs are 0-index, so we don't have to do any math when calling their referenced arrays
	ctx1 := context.Background()
	ctx1 = context.WithValue(ctx1, "serverId", "0")
	ctx2 := context.Background()
	ctx2 = context.WithValue(ctx2, "serverId", "1")

	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err.Error())
		return
	}

	firstTime := db.IsFirstRun()

	if firstTime {
		err := generate.Generate(25)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx1 = context.WithValue(ctx1, "db", db)
	ctx2 = context.WithValue(ctx2, "db", db)

	wsServers[0] = &WebsocketServer{0, util.RandomGUID(), fmt.Sprintf("wss://localhost:%v/eventsub", port), []WebsocketConnection{}, false}
	wsServers[1] = &WebsocketServer{1, util.RandomGUID(), fmt.Sprintf("wss://localhost:%v/eventsub", port+1), []WebsocketConnection{}, true}

	RegisterHandlers(m)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	s1 := StartIndividualServer(port, reconnectTestTimer, m, ctx1)
	s2 := StartIndividualServer(port+1, 0, m, ctx2) // Start second server, at a port above. Never has a reconnect timer

	<-stop // Wait for ctrl + c

	log.Print("Shutting down EventSub WebSocket server ...\n")
	db.DB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	if err1 := s1.Shutdown(ctx); err1 != nil {
		log.Fatal(err1)
	}

	if err2 := s2.Shutdown(ctx); err2 != nil {
		log.Fatal(err2)
	}
}

func StartIndividualServer(port int, reconnectTestTimer int, m *http.ServeMux, ctx context.Context) http.Server {
	s := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: m,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Mock EventSub websocket server started on port %d", port)

		go func() {
			if reconnectTestTimer != 0 {
				log.Printf("Reconnect testing enabled. Will be sent in %d seconds.", reconnectTestTimer)

				select {
				case <-time.After(time.Second * time.Duration(reconnectTestTimer)):
					activateReconnectTest(s, ctx)
				}
			}
		}()

		if err := s.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	return s
}

func RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("/debug", debugPageHandler)
	m.HandleFunc("/eventsub", eventsubHandle)
}

func debugPageHandler(w http.ResponseWriter, r *http.Request) {
	debugTemplate.Execute(w, DebugServerButtons{Server1: wsServers[0].connectionUrl, Server2: wsServers[1].connectionUrl})
}

// Used for template weirdness
type DebugServerButtons struct {
	Server1 string
	Server2 string
}

var debugTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }

		let wsUrl = "{{.Server1}}";
		if (evt.altKey || evt.ctrlKey) {
			wsUrl = "{{.Server2}}"
		}

        ws = new WebSocket(wsUrl);
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RECEIVED: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
            console.error("ERROR", evt);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server ({{.Server1}}), 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<p>For testing the higher port ({{.Server2}}), hold ctrl OR alt and click "Open"</p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
