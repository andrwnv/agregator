package main

import (
	"comment-service/handlers"
	"comment-service/websockets"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router        *mux.Router
	WebsocketPool *websockets.Pool
}

func (s *Server) Init() {
	s.Router = mux.NewRouter()
	s.WebsocketPool = websockets.NewPool()

	// s.Router.PathPrefix("/comments")
	s.setupRoutes()
}

func (s *Server) GetRedirect(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.Router.HandleFunc(url, handler).Methods("GET")
}

func (s *Server) PostRedirect(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.Router.HandleFunc(url, handler).Methods("POST")
}

func (s *Server) DeleteRedirect(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.Router.HandleFunc(url, handler).Methods("DELETE")
}

func (s *Server) PatchRedirect(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.Router.HandleFunc(url, handler).Methods("PATCH")
}

func (s *Server) setupRoutes() {
	s.GetRedirect("/test", handlers.Hello)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(s.WebsocketPool, w, r)
	})
}

// ----------------------

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

func main() {
	server := &Server{}
	server.Init()

	go server.WebsocketPool.Start()

	fmt.Println("Test websocket app")
	log.Fatal(http.ListenAndServe(":3030", server.Router))
}
