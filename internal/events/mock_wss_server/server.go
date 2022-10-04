package mock_wss_server

import (
	"context"
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

var upgrader = websocket.Upgrader{}
var connections []WebsocketConnection

var debug = false

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	// Connection sucessful. WebSocket is now open.

	// RFC3339Nano = "2022-10-04T12:38:15.479Z" ; Specifically, it includes milliseconds, whereas RFC3339 doesn't.
	connectedAtTimestamp := time.Now().UTC().Format(time.RFC3339Nano)

	// Generate unique websocket ID.
	websocketId := util.RandomGUID()

	// Add to websocket connection list.
	connections = append(connections, WebsocketConnection{ID: websocketId})
	printConnections()

	// Send "websocket_welcome" message
	welcomeMsg := []byte(fmt.Sprintf(`
{
    "metadata": {
        "message_id": "` + util.RandomGUID() + `",
        "message_type": "websocket_welcome",
        "message_timestamp": "` + time.Now().UTC().Format(time.RFC3339Nano) + `"
    },
    "payload": {
        "websocket": {
            "id": "` + websocketId + `",
            "status": "connected",
            "minimum_message_frequency_seconds": ` + fmt.Sprintf("%d", MINIMUM_MESSAGE_FREQUENCY_SECONDS) + `,
            "connected_at": "` + connectedAtTimestamp + `"
        }
    }
}
`)) // TODO: Convert to proper JSON in Go
	conn.WriteMessage(1, welcomeMsg)

	if debug {
		log.Printf("[DEBUG] Write: %s", welcomeMsg)
	}

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			close(websocketId)
			break
		}
		log.Printf("recv: [%d] %s", mt, message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			close(websocketId)
			break
		}
	}
}

func close(websocketId string) {
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

func StartServer(port int, enableDebug bool) {
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
	m.HandleFunc("/", home)
	m.HandleFunc("/echo", echo)
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "wss://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
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
            print("RESPONSE: " + evt.data);
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
	ID string
}
