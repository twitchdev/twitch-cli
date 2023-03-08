package mock_ws_server

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	clientId                 string
	conn                     *websocket.Conn
	mutex                    sync.Mutex
	clientConnectedTimestamp string

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
