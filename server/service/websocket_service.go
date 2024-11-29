package service

import (
	"sync"

	"github.com/gorilla/websocket"
)

type IWebsocketConnection interface {
	SendMessage(message []byte) error
	ReceiveMessage() ([]byte, error)
	Close() error
}

type WebSocketConnection struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

func (ws *WebSocketConnection) SendMessage(message []byte) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return ws.conn.WriteMessage(websocket.TextMessage, message)
}

func (ws *WebSocketConnection) ReceiveMessage() ([]byte, error) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	_, msg, err := ws.conn.ReadMessage()
	return msg, err
}

func (ws *WebSocketConnection) Close() error {
	return ws.conn.Close()
}

func NewWebSocketConnection(conn *websocket.Conn) *WebSocketConnection {
	return &WebSocketConnection{
		conn: conn,
	}
}
