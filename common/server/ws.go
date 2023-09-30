package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/interactivehub/engine/common/ws"
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

type WSInteractor struct {
	WSListener
	WSWriter
}

type WSListener interface {
	ListenEvents()
}

type WSWriter interface {
	WriteEvent(t string, p interface{}) error
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

func RunWSServer() {
	port := os.Getenv("WS_PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	RunWSServerOnAddr(addr)
}

func RunWSServerOnAddr(addr string) {
	wsServer := NewWSServer()

	http.HandleFunc("/ws", wsServer.ServeHTTP)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func (s *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	s.client = conn

	defer s.client.Close()

	s.ListenEvents()
}

func (s *WSServer) ListenEvents() {
	for {
		var event ws.Event

		err := s.client.ReadJSON(&event)
		if err != nil {
			log.Println(err)
		}

		log.Println(event)
	}
}

func (s *WSServer) WriteEvent(t string, p interface{}) error {
	if s.client == nil {
		return errors.New("missing ws client")
	}

	event := ws.NewEvent(t, p)

	return s.client.WriteJSON(event)
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
