package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mohammadzaidhussain/pizza-shop/logger"
	"github.com/mohammadzaidhussain/pizza-shop/service"
)

type IWebSocketHandler interface {
	HandleConnection(ctx *gin.Context)
	GetConnectionMap() *map[string]service.IWebsocketConnection
}

type WebSocketHandler struct {
	upgrader   websocket.Upgrader
	connection *map[string]service.IWebsocketConnection
	mutex      sync.Mutex
}

func (h *WebSocketHandler) HandleConnection(ctx *gin.Context) {
	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to upgrade connection"))
		return
	}
	defer conn.Close()

	conn.WriteMessage(websocket.TextMessage, []byte("Started taking order"))

	connection := service.NewWebSocketConnection(conn)
	h.addConnection("pizza", connection)

	for {
		logger.Log("no message is coming from client")
	}
}

func (h *WebSocketHandler) addConnection(clientId string, connection service.IWebsocketConnection) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	(*h.connection)[clientId] = connection
}

func (h *WebSocketHandler) GetConnectionMap() *map[string]service.IWebsocketConnection {
	return h.connection
}

func GetNewWebSocketHandler() *WebSocketHandler {
	connection := make(map[string]service.IWebsocketConnection)
	return &WebSocketHandler{
		connection: &connection,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}
