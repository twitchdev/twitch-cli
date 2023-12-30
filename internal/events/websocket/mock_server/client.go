package mock_server

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	clientName           string // Unique name for the client. Not the Client ID.
	conn                 *websocket.Conn
	mutex                sync.Mutex
	ConnectedAtTimestamp string // RFC3339Nano timestamp indicating when the client connected to the server
	connectionUrl        string
	KeepAliveEnabled     bool

	mustSubscribeTimer *time.Timer
	keepAliveChanOpen  bool
	keepAliveLoopChan  chan struct{}
	keepAliveTimer     *time.Ticker
	pingChanOpen       bool
	pingLoopChan       chan struct{}
	pingTimer          *time.Ticker
}

func (c *Client) SendMessage(messageType int, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.conn.WriteMessage(messageType, data)
}

func (c *Client) CloseWithReason(reason *CloseMessage) {
	c.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(reason.code, reason.message),
		time.Now().Add(2*time.Second),
	)
}

func (c *Client) CloseDirty() {
	c.conn.Close()
}
