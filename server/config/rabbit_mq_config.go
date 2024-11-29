package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mohammadzaidhussain/pizza-shop/logger"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConection struct {
	conn  *amqp091.Connection
	queue string
}

func GetNewRabbitMQConnection() *RabbitMQConection {
	host := GetEnvProperty("rabbit_mq_host")
	port := GetEnvProperty("rabbit_mq_port")
	username := GetEnvProperty("rabbit_mq_username")
	password := GetEnvProperty("rabbit_mq_password")
	queue := GetEnvProperty("rabbit_mq_default_queue")

	PORT, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("invalid port : %v", err))
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, PORT)

	fmt.Println(url)

	conn, err := amqp091.Dial(url)

	if err != nil {
		panic(fmt.Sprintf("failed to connect to RabbitMQ : %v", err))
	}

	log.Println("RabbitMq has been connected")

	return &RabbitMQConection{
		conn:  conn,
		queue: queue,
	}
}

func (r *RabbitMQConection) Connect() *amqp091.Connection {
	host := GetEnvProperty("rabbit_mq_host")
	port := GetEnvProperty("rabbit_mq_port")
	username := GetEnvProperty("rabbit_mq_username")
	password := GetEnvProperty("rabbit_mq_password")

	PORT, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("invalid port : %v", err))
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, PORT)

	conn, err := amqp091.Dial(url)

	if err != nil {
		panic(fmt.Sprintf("failed to connect to RabbitMQ : %v", err))
	}

	log.Println("RabbitMq has been reconnected")
	return conn
}

func (r *RabbitMQConection) DeclareQueue(queueName string) error {
	var err error
	channel, err := r.conn.Channel()
	defer channel.Close()

	_, err = channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (r *RabbitMQConection) GetConnection() *amqp091.Connection {
	if r.conn == nil {
		r.conn = r.Connect()
	}
	return r.conn
}

func (r *RabbitMQConection) GetChannel() *amqp091.Channel {
	var channel *amqp091.Channel
	connection := r.conn
	if connection == nil {
		connection = r.Connect()
	}
	channel, err := r.conn.Channel()
	if err != nil {
		channel, _ = r.conn.Channel()
		logger.Log("channel is nil")
	}
	if channel != nil && channel.IsClosed() {
		channel, err = r.conn.Channel()
		if err != nil {
			logger.Log("channel was closed and error creating channel")
		}
	}
	return channel
}

func (r *RabbitMQConection) GetQueue() string {
	return r.queue
}

func (r *RabbitMQConection) Close() {
	r.conn.Close()
}
