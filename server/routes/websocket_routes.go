package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohammadzaidhussain/pizza-shop/handler"
)

func RegisterWebSocketRoutes(router *gin.RouterGroup, websocketHandler handler.IWebSocketHandler) {
	router.GET(
		"/",
		websocketHandler.HandleConnection,
	)
}
