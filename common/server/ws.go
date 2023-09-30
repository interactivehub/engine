package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}
)

type WSServer struct {
	upgrader websocket.Upgrader
	client   *websocket.Conn
}

func NewWSServer() *WSServer {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}

	return &WSServer{
		upgrader: upgrader,
		client:   nil,
	}
}

func RunWSServer(f func(server *WSServer)) {
	port := os.Getenv("WS_PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	RunWSServerOnAddr(addr, f)
}

func RunWSServerOnAddr(addr string, f func(server *WSServer)) {
	log.Printf("Starting WebSocket server on port %s", addr)

	wsServer := NewWSServer()

	http.HandleFunc("/ws", wsServer.ServeHTTP(f))

	log.Fatal(http.ListenAndServe(addr, nil))
}

func (s *WSServer) ServeHTTP(f func(server *WSServer)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		s.client = conn

		f(s)

		defer s.client.Close()
	}
}

func (s *WSServer) Client() *websocket.Conn {
	return s.client
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	hubURL := os.Getenv("HUB_URL")
	if hubURL == "" {
		hubURL = "https://localhost:5173"
	}

	switch origin {
	case hubURL:
		return true
	default:
		return false
	}
}
