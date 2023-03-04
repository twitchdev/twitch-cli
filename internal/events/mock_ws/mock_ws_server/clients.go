package mock_ws_server

type Clients struct {
	Clients map[string]*Client
}

func (c *Clients) Get(key string) (*Client, bool) {
	client, ok := c.Clients[key]
	return client, ok
}

func (c *Clients) Put(key string, client *Client) {
	c.Clients[key] = client
}

func (c *Clients) Delete(key string) {
	delete(c.Clients, key)
}

func (c *Clients) Length() int {
	n := len(c.Clients)
	return n
}

func (c *Clients) All() []*Client {
	clients := []*Client{}
	for _, client := range c.Clients {
		clients = append(clients, client)
	}
	return clients
}
