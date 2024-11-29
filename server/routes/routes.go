package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohammadzaidhussain/pizza-shop/handler"
	"github.com/mohammadzaidhussain/pizza-shop/service"
)

func RegisterRoutes(r *gin.Engine, messagePublisher service.IMessagePublisher, websocketHandler handler.IWebSocketHandler) {

	router := r.Group("/")

	wsr := router.Group("/ws")
	{
		RegisterWebSocketRoutes(wsr, websocketHandler)
	}

	or := router.Group("/orders")
	{
		RegisterOrderRoutes(or, messagePublisher)
	}

}
