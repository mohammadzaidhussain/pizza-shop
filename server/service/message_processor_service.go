package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/mohammadzaidhussain/pizza-shop/constants"
	"github.com/mohammadzaidhussain/pizza-shop/logger"
	"github.com/mohammadzaidhussain/pizza-shop/utils"
	"github.com/rabbitmq/amqp091-go"
)

type IMessageProcessor interface {
	ProcessMessage(message interface{}) error
}

type MessageProcessor struct {
	publisher  IMessagePublisher
	connection *map[string]IWebsocketConnection
	mutex      sync.RWMutex
}

func (mp *MessageProcessor) ProcessMessage(message interface{}) error {
	msg := message.(amqp091.Delivery)
	var event map[string]interface{}
	var err error
	if err = json.Unmarshal(msg.Body, &event); err != nil {
		logger.Log(fmt.Sprintf("Failed to unmarshal message : %v", err))
		msg.Nack(false, true)
		return err
	}

	logger.Log(fmt.Sprintf("processing message: %v", event))

	if val, ok := event["order_status"]; ok {
		switch val {
		case constants.ORDER_ORDERED:
			{
				err = mp.handleOrderOrdered(event)
			}
		case constants.ORDER_PREPARING:
			{
				err = mp.handleOrderPreparing(event)
			}
		case constants.ORDER_PREPARED:
			{
				err = mp.handleOrderPrepared(event)
			}
		default:
			{
				logger.Log("No order to be processed!")
			}
		}
		if err != nil {
			logger.Log(fmt.Sprintf("Error Processing Message: %v", err))
			msg.Nack(false, true)
			return err
		}
	}

	msg.Ack(false)
	return nil
}

func (mp *MessageProcessor) handleOrderOrdered(event map[string]interface{}) error {
	var err error
	logger.Log(fmt.Sprintf("order %v accepted", event))
	event["order_status"] = constants.ORDER_PREPARING
	err = mp.publisher.PublishEvent(constants.KITCHEN_ORDER_QUEUE, event)
	if err != nil {
		logger.Log(fmt.Sprintf("error: %v, event: %v", err, event))
		message := map[string]string{
			"message": constants.ORDER_CANCELLED,
			"error":   err.Error(),
		}
		messagesBytes, _ := json.Marshal(message)
		if mp.connection != nil {
			mp.mutex.Lock()
			defer mp.mutex.Unlock()
			pizza := (*mp.connection)["pizza"]
			if pizza != nil {
				err = (*mp.connection)["pizza"].SendMessage(messagesBytes)
			}
		}
	}
	return err
}

func (mp *MessageProcessor) handleOrderPreparing(event map[string]interface{}) error {
	var err error
	logger.Log(fmt.Sprintf("order %v accepted", event))
	event["order_status"] = constants.ORDER_PREPARED
	time.Sleep(utils.GenerateRandomDuration(1, 6))
	err = mp.publisher.PublishEvent(constants.KITCHEN_ORDER_QUEUE, event)
	if err != nil {
		logger.Log(fmt.Sprintf("error: %v, event: %v", err, event))
		message := map[string]string{
			"message": constants.ORDER_CANCELLED,
			"error":   err.Error(),
		}
		messagesBytes, _ := json.Marshal(message)
		if mp.connection != nil {
			mp.mutex.Lock()
			defer mp.mutex.Unlock()
			pizza := (*mp.connection)["pizza"]
			if pizza != nil {
				err = (*mp.connection)["pizza"].SendMessage(messagesBytes)
			}
		}
	}
	return err
}

func (mp *MessageProcessor) handleOrderPrepared(event map[string]interface{}) error {
	var err error
	logger.Log(fmt.Sprintf("order %v prepared successfully", event["order_no"]))
	event["order_status"] = constants.ORDER_DELIVERED
	logger.Log(fmt.Sprintf("error: %v, event: %v", err, event))
	message := map[string]interface{}{
		"message": constants.ORDER_PREPARED_SUCCESSFULLY,
		"order":   event,
	}
	messagesBytes, _ := json.Marshal(message)
	if mp.connection != nil {
		mp.mutex.Lock()
		defer mp.mutex.Unlock()
		pizza := (*mp.connection)["pizza"]
		if pizza != nil {
			err = (*mp.connection)["pizza"].SendMessage(messagesBytes)
		}
	}
	return err
}

func GetMessageProcessorService(publisher IMessagePublisher, connection *map[string]IWebsocketConnection) *MessageProcessor {
	return &MessageProcessor{
		publisher:  publisher,
		connection: connection,
	}
}
