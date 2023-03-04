package mock_ws_server

type Servers struct {
	Servers map[string]*WebSocketServer
}

func (c *Servers) Get(key string) (*WebSocketServer, bool) {
	server, ok := c.Servers[key]
	return server, ok
}

func (c *Servers) Put(key string, server *WebSocketServer) {
	c.Servers[key] = server
}

func (c *Servers) Delete(key string) {
	delete(c.Servers, key)
}

func (c *Servers) Length() int {
	n := len(c.Servers)
	return n
}

func (c *Servers) All() []*WebSocketServer {
	servers := []*WebSocketServer{}
	for _, server := range c.Servers {
		servers = append(servers, server)
	}
	return servers
}
