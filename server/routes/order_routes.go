package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohammadzaidhussain/pizza-shop/handler"
	"github.com/mohammadzaidhussain/pizza-shop/service"
)

func RegisterOrderRoutes(router *gin.RouterGroup, messagePublisher service.IMessagePublisher) {

	oh := handler.GetOrderHandler(messagePublisher)

	router.POST(
		"/create",
		oh.CreateOrder,
	)
}
