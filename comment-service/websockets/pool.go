package websockets

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Current pool size = ", len(pool.Clients))

			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User joined"})
			}
			break

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Current pool size = ", len(pool.Clients))

			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User disconnect"})
			}
			break

		case message := <-pool.Broadcast:
			fmt.Println("Send msg for all")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
