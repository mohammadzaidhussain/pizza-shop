package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mohammadzaidhussain/pizza-shop/config"
	"github.com/mohammadzaidhussain/pizza-shop/logger"
	"github.com/rabbitmq/amqp091-go"
)

type IMessagePublisher interface {
	DeclareQueue(queueName string) error
	PublishEvent(queueName string, body interface{}) error
}

type MessagePublisher struct {
	conf *config.RabbitMQConection
}

func (mp *MessagePublisher) DeclareQueue(queueName string) error {
	channel := mp.conf.GetChannel()
	if channel == nil {
		return fmt.Errorf("Message channel is nil, please retry!")
	}

	_, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (mp *MessagePublisher) PublishEvent(queueName string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if queueName == "" {
		queueName = config.GetEnvProperty("rabbit_mq_default_queue")
	}

	channel := mp.conf.GetChannel()
	if channel == nil {
		panic("Messaging channel is nil, retry !")
	}
	if channel.IsClosed() {
		panic("could not publish event, channel closed")
	}

	logger.Log(fmt.Sprintf("created new channel....%v", &channel))

	err = channel.PublishWithContext(ctx,
		"",
		queueName,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp091.Persistent,
		},
	)

	if err != nil {
		return err
	}

	logger.Log(fmt.Sprintf("Event published: %v", body))
	channel.Close()
	logger.Log(fmt.Sprintf("channel closed: %v", &channel))

	return nil

}

func GetMessagePublisherService() *MessagePublisher {
	rabbitMQConf := config.GetNewRabbitMQConnection()
	return &MessagePublisher{
		conf: rabbitMQConf,
	}
}
