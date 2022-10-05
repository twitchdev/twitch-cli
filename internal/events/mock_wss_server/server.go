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
	"text/template"
	"time"

	"github.com/gorilla/websocket"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/util"
)

/* Minimum time between messages before the server disconnects a client. AKA, "timeout period"
 * This is a const for now, but it may need to be a flag in the future. */
const MINIMUM_MESSAGE_FREQUENCY_SECONDS = 10

var connectionURL string

var upgrader = websocket.Upgrader{}
var connections []WebsocketConnection
var websocketId string
var connectedAtTimestamp string
var shuttingDown = false // Used during reconnect to avoid accepting new messages

var debug = false

func eventsubHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	// Connection sucessful. WebSocket is now open.

	// RFC3339Nano = "2022-10-04T12:38:15.548912638Z07" ; This is used by Twitch in production
	connectedAtTimestamp = time.Now().UTC().Format(time.RFC3339Nano)
	conn.SetReadDeadline(time.Now().Add(time.Second * MINIMUM_MESSAGE_FREQUENCY_SECONDS))

	// Generate unique websocket ID.
	websocketId = util.RandomGUID()

	// Add to websocket connection list.
	connections = append(connections, WebsocketConnection{ID: websocketId, Conn: conn})
	printConnections()

	// TODO: Decline websocket if it reached 100 connections from the same application access token
	// Use "shuttingDown" bool for this

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
					ID:                             websocketId,
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
			onCloseConnection(websocketId)
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

func onCloseConnection(websocketId string) {
	log.Printf("Disconnected websocket %s", websocketId)

	// Remove from list
	c := 0
	for i := 0; i < len(connections); i++ {
		if connections[i].ID == websocketId {
			c = i
			break
		}
	}
	connections = append(connections[:c], connections[c+1:]...)

	printConnections()
}

func printConnections() {
	currentConnections := ""
	for _, s := range connections {
		currentConnections += s.ID + ", "
	}
	if currentConnections != "" {
		currentConnections = string(currentConnections[:len(currentConnections)-2])
	}
	log.Printf("Connections: (%d) [ %s ]", len(connections), currentConnections)
}

func terminateServer(server http.Server, ctx context.Context, reconnect bool) {
	timer := 0
	if reconnect {
		timer = 30
		log.Println("Terminating server with reconnect notice. Shutting down in 30 seconds.")
	} else {
		log.Println("Terminating server immediately.")
	}

	if reconnect {
		// Stop processing new messages
		shuttingDown = true

		// Send reconnect notices
		for _, c := range connections {
			reconnectMsg, _ := json.Marshal(
				ReconnectMessage{
					Metadata: MessageMetadata{
						MessageID:        util.RandomGUID(),
						MessageType:      "websocket_reconnect",
						MessageTimestamp: time.Now().UTC().Format(time.RFC3339Nano),
					},
					Payload: ReconnectMessagePayload{
						Websocket: ReconnectMessagePayloadWebsocket{
							ID:                             websocketId,
							Status:                         "reconnecting",
							MinimumMessageFrequencySeconds: MINIMUM_MESSAGE_FREQUENCY_SECONDS,
							Url:                            connectionURL,
							ConnectedAt:                    connectedAtTimestamp,
							ReconnectingAt:                 time.Now().UTC().Add(time.Second * 30).Format(time.RFC3339Nano),
						},
					},
				},
			)

			c.Conn.WriteMessage(1, reconnectMsg) // Ignore err here because we're shutting down so it doesn't really matter
		}
	}

	// Wait 30 seconds if reconnect, otherwise it'll execute immediately (0 second timer)
	select {
	case <-time.After(time.Second * time.Duration(timer)):
		for _, c := range connections {
			c.Conn.Close()
		}

		server.Shutdown(ctx)
	}
}

func StartServer(port int, enableDebug bool, reconnectTestTimer int) {
	debug = enableDebug

	m := http.NewServeMux()

	ctx := context.Background()

	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err.Error())
		return
	}

	firstTime := db.IsFirstRun()

	if firstTime {
		//err := generate.Generate(25)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}

	ctx = context.WithValue(ctx, "db", db)

	RegisterHandlers(m)
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
		log.Print("Mock EventSub websocket server started")

		go func() {
			if reconnectTestTimer != 0 {
				log.Printf("Reconnect testing enabled. Will be sent in %d seconds.", reconnectTestTimer)

				select {
				case <-time.After(time.Second * time.Duration(reconnectTestTimer)):
					terminateServer(s, ctx, true)
				}
			}
		}()

		if err := s.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-stop

	log.Print("Shutting down EventSub WebSocket server ...\n")
	db.DB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("/debug", debugPageHandler)
	m.HandleFunc("/eventsub", eventsubHandle)
}

func debugPageHandler(w http.ResponseWriter, r *http.Request) {
	// Set connection URL for future use
	// TODO: Find an earlier spot for this. If this debug page is ever removed, then reconnection is broken due to lack of URL
	connectionURL = "wss://" + r.Host + "/eventsub"

	debugTemplate.Execute(w, connectionURL)
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
        ws = new WebSocket("{{.}}");
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
            console.error(evt);
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
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
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

type WebsocketConnection struct {
	ID   string
	Conn *websocket.Conn
}
