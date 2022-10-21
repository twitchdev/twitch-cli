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
	"path/filepath"
	"strconv"
	"sync"
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

var upgrader = websocket.Upgrader{}
var debug = false

// List of websocket servers. Limited to 2 for now, as allowing more would require rewriting reconnect stuff.
var wsServers = [2]*WebsocketServer{}

type WebsocketServer struct {
	serverId             int                    // Int representing the ID of the server (0, 1, 2, ...)
	websocketId          string                 // UUID of the websocket. Used for subscribing via EventSub
	connectionUrl        string                 // URL used to connect to the websocket. Used for reconnect messages
	connections          []*WebsocketConnection // Current clients connected to this websocket
	deactivatedStatus    bool                   // Boolean used for preventing connections/messages during deactivation; Used for reconnect testing
	reconnectTestTimeout int                    // Timeout for reconnect testing after first client connects; 0 if reconnect testing not enabled.
	firstClientConnected bool                   // Whether or not the first client has connected (used for reconnect testing)
}

type WebsocketConnection struct {
	clientId             string
	Conn                 *websocket.Conn
	mu                   sync.Mutex
	connectedAtTimestamp string
	pingLoopChan         chan struct{}
	kaLoopChan           chan struct{}
	closed               bool
}

func (wc *WebsocketConnection) SendMessage(messageType int, data []byte) error {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	return wc.Conn.WriteMessage(messageType, data)
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

	// Activate reconnect testing upon first client connection
	if !wsSrv.firstClientConnected {
		wsSrv.firstClientConnected = true

		if wsSrv.reconnectTestTimeout != 0 {
			go func() {
				log.Printf("First client connected; Reconnect testing enabled. Notices will be sent in %d seconds.", wsSrv.reconnectTestTimeout)

				select {
				case <-time.After(time.Second * time.Duration(wsSrv.reconnectTestTimeout)):
					activateReconnectTest(r.Context())
				}
			}()
		}
	}

	// TODO: Decline websocket if it reached 100 connections from the same application access token

	// RFC3339Nano = "2022-10-04T12:38:15.548912638Z07" ; This is used by Twitch in production
	connectedAtTimestamp := time.Now().UTC().Format(time.RFC3339Nano)
	conn.SetReadDeadline(time.Now().Add(time.Second * MINIMUM_MESSAGE_FREQUENCY_SECONDS))

	// Add to websocket connection list.
	wc := &WebsocketConnection{
		clientId:             util.RandomGUID(),
		Conn:                 conn,
		connectedAtTimestamp: connectedAtTimestamp,
		closed:               false,
	}
	wsSrv.connections = append(wsSrv.connections, wc)
	printConnections(wsSrv.serverId)

	// Send "websocket_welcome" message
	welcomeMsg, _ := json.Marshal(
		WelcomeMessage{
			Metadata: MessageMetadata{
				MessageID:        util.RandomGUID(),
				MessageType:      "session_welcome",
				MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
			},
			Payload: WelcomeMessagePayload{
				Session: WelcomeMessagePayloadSession{
					ID:                             wsSrv.websocketId,
					Status:                         "connected",
					MinimumMessageFrequencySeconds: MINIMUM_MESSAGE_FREQUENCY_SECONDS,
					ReconnectUrl:                   nil,
					ConnectedAt:                    connectedAtTimestamp,
				},
			},
		},
	)
	wc.SendMessage(websocket.TextMessage, welcomeMsg)
	if debug {
		log.Printf("[DEBUG] Write: %s", welcomeMsg)
	}

	// TODO: Look to implement a way to shut off pings. This would be used specifically for testing the timeout feature.
	// Set up ping/pong handling
	pingTicker := time.NewTicker(5 * time.Second)
	wc.pingLoopChan = make(chan struct{}) // Also used for keepalive
	go func() {
		// Set pong handler
		// Weirdly, pongs are not seen as messages read by conn.ReadMessage, so we have to reset the deadline manually
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(time.Second * MINIMUM_MESSAGE_FREQUENCY_SECONDS))
			return nil
		})

		// Ping loop
		for {
			select {
			case <-wc.pingLoopChan:
				pingTicker.Stop()
				return
			case <-pingTicker.C:
				err := wc.SendMessage(websocket.PingMessage, []byte{})
				if err != nil {
					onCloseConnection(wsSrv.serverId, wc)
				}
			}
		}
	}()

	// Set up keepalive loop
	kaTicker := time.NewTicker(10 * time.Second)
	wc.kaLoopChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-wc.kaLoopChan:
				kaTicker.Stop()
			case <-kaTicker.C:
				keepAliveMsg, _ := json.Marshal(
					KeepaliveMessage{
						Metadata: MessageMetadata{
							MessageID:        util.RandomGUID(),
							MessageType:      "session_keepalive",
							MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
						},
						Payload: KeepaliveMessagePayload{},
					},
				)
				err := wc.SendMessage(websocket.TextMessage, keepAliveMsg)
				if err != nil {
					onCloseConnection(wsSrv.serverId, wc)
				}

				if debug {
					log.Printf("[DEBUG] Write: %s", keepAliveMsg)
				}
			}
		}
	}()

	// TODO: Read messages
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * MINIMUM_MESSAGE_FREQUENCY_SECONDS))
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			onCloseConnection(wsSrv.serverId, wc)
			break
		}
		if debug {
			log.Printf("recv: [%d] %s", mt, message)
		}
		/*err = wc.SendMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			onCloseConnection(websocketId)
			break
		}*/
	}
}

func onCloseConnection(serverId int, wc *WebsocketConnection) {
	// Close ping loop chan
	if !wc.closed {
		close(wc.pingLoopChan)
		close(wc.kaLoopChan)
	}
	wc.closed = true

	wsSrv := wsServers[serverId]

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

	printConnections(wsSrv.serverId)
}

func printConnections(serverId int) {
	currentConnections := ""
	wsSrv := wsServers[serverId]
	for _, s := range wsSrv.connections {
		currentConnections += s.clientId + ", "
	}
	if currentConnections != "" {
		currentConnections = string(currentConnections[:len(currentConnections)-2])
	}
	log.Printf("[Server %v] Connections: (%d) [ %s ]", serverId, len(wsSrv.connections), currentConnections)
}

func activateReconnectTest(ctx context.Context) {
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
	log.Printf("Server \"not accepting connections\" status: [Server 0: %v, Server 1: %v]", wsServers[0].deactivatedStatus, wsServers[1].deactivatedStatus)

	if debug {
		log.Printf("Connections at time of close: %v", len(wsSrv.connections))
	}

	// Send reconnect notices
	for _, c := range wsSrv.connections {
		reconnectMsg, _ := json.Marshal(
			ReconnectMessage{
				Metadata: MessageMetadata{
					MessageID:        util.RandomGUID(),
					MessageType:      "session_reconnect",
					MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
				},
				Payload: ReconnectMessagePayload{
					Session: ReconnectMessagePayloadSession{
						ID:                             wsSrv.websocketId,
						Status:                         "reconnecting",
						MinimumMessageFrequencySeconds: nil,
						ReconnectUrl:                   wsSrv.connectionUrl,
						ConnectedAt:                    c.connectedAtTimestamp,
					},
				},
			},
		)

		err := c.SendMessage(websocket.TextMessage, reconnectMsg)
		if err != nil {
			log.Printf("ERROR (clientId %v): %v", c.clientId, err)
		} else {
			log.Printf("Sent reconnect notice to %v", c.clientId)
		}
	}

	log.Printf("Reconnect notices sent for server %v", serverId)
	log.Printf("Use this new URL for connections: %v", wsAltSrv.connectionUrl)

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

func StartServer(port int, enableDebug bool, reconnectTestTimer int, sslEnabled bool) {
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

	wsSrv1Url := fmt.Sprintf("ws://localhost:%v/eventsub", port)
	wsSrv2Url := fmt.Sprintf("ws://localhost:%v/eventsub", port+1)

	// Change to WSS if SSL is enabled via flag
	if sslEnabled {
		wsSrv1Url = fmt.Sprintf("wss://localhost:%v/eventsub", port)
		wsSrv2Url = fmt.Sprintf("wss://localhost:%v/eventsub", port+1)
	}

	wsServers[0] = &WebsocketServer{
		serverId:             0,
		websocketId:          util.RandomGUID(),
		connectionUrl:        wsSrv1Url,
		connections:          []*WebsocketConnection{},
		deactivatedStatus:    false,
		reconnectTestTimeout: reconnectTestTimer,
		firstClientConnected: false,
	}
	wsServers[1] = &WebsocketServer{
		serverId:             1,
		websocketId:          util.RandomGUID(),
		connectionUrl:        wsSrv2Url,
		connections:          []*WebsocketConnection{},
		deactivatedStatus:    true, // 2nd server is deactivated by default. Will reactivate for reconnect testing.
		reconnectTestTimeout: 0,    // No reconnect testing
		firstClientConnected: false,
	}

	RegisterHandlers(m)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	s1 := StartIndividualServer(port, reconnectTestTimer, sslEnabled, m, ctx1)
	s2 := StartIndividualServer(port+1, 0, sslEnabled, m, ctx2) // Start second server, at a port above. Never has a reconnect timer

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

func StartIndividualServer(port int, reconnectTestTimer int, sslEnabled bool, m *http.ServeMux, ctx context.Context) http.Server {
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

		if sslEnabled { // Open HTTP server with HTTPS support
			home, _ := util.GetApplicationDir()
			crtFile := filepath.Join(home, "localhost.crt")
			keyFile := filepath.Join(home, "localhost.key")

			if err := s.ListenAndServeTLS(crtFile, keyFile); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalf(`%v
	** You need to generate localhost.crt and localhost.key for this to work **
	** Please run these commands (Note: you'll have a cert error in your web browser, but it'll still start): **
		openssl genrsa -out "%v" 2048
		openssl req -new -x509 -sha256 -key "%v" -out "%v" -days 3650`,
						err, keyFile, keyFile, crtFile)
				}
			}
		} else { // Open HTTP server without HTTPS support
			if err := s.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalf("%v", err)
				}
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
            print("OPEN - " + ws.url);
        }
        ws.onclose = function(evt) {
            print("CLOSE - " + ws.url);
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("[RECEIVED / " + new Date().toISOString() + "]: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("[ERROR / " + new Date().toISOString() + "]: " + evt.data);
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
	document.getElementById("clear").onclick = function(evt) {
		output.innerHTML = "";
	}
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
<button id="send">Send</button><br><br>
<button id="clear">Clear</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
