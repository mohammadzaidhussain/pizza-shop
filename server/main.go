package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mohammadzaidhussain/pizza-shop/config"
	"github.com/mohammadzaidhussain/pizza-shop/constants"
	"github.com/mohammadzaidhussain/pizza-shop/handler"
	"github.com/mohammadzaidhussain/pizza-shop/logger"
	"github.com/mohammadzaidhussain/pizza-shop/middleware"
	"github.com/mohammadzaidhussain/pizza-shop/routes"
	"github.com/mohammadzaidhussain/pizza-shop/service"
)

func main() {

	app := gin.Default()

	app.Use(gin.Recovery())

	app.Use(middleware.CorsMiddleware)

	app.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":    "Pizza Shop is open",
			"statusCode": 200,
		})
	})

	messagePublisher := service.GetMessagePublisherService()
	messageConsumer := service.GetMessageConsumerService()

	websocketHandler := handler.GetNewWebSocketHandler()
	messageProcessor := service.GetMessageProcessorService(messagePublisher, websocketHandler.GetConnectionMap())

	go func() {
		err := messageConsumer.ConsumeEventAndProcess(constants.KITCHEN_ORDER_QUEUE, messageProcessor)
		if err != nil {
			logger.Log(fmt.Sprintf("failed to consume events : %v", err))
		}
	}()

	routes.RegisterRoutes(app, messagePublisher, websocketHandler)

	port := config.GetEnvProperty("port")
	logger.Log(fmt.Sprintf("Pizza shop started successfully on port : %s", port))

	app.Run(fmt.Sprintf(":%s", port))

}
