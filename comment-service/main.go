package main

import (
	"comment-service/websockets"
	"fmt"
	"net/http"
)

func serveWebsocket(pool *websockets.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket Endpoint")
	conn, err := websockets.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websockets.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websockets.NewPool()
	go pool.Start() // GORUTINE LMAO

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(pool, w, r)
	})
}

func main() {
	fmt.Println("Test websocket app")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
