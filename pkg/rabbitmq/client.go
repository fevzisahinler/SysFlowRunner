package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	Connection  *amqp.Connection
	rabbitmqURL string
}

func NewRabbitMQClient(rabbitmqURL string) (*RabbitMQClient, error) {
	conn, err := connectRabbitMQ(rabbitmqURL)
	if err != nil {
		return nil, err
	}
	return &RabbitMQClient{
		Connection:  conn,
		rabbitmqURL: rabbitmqURL,
	}, nil
}

func connectRabbitMQ(rabbitmqURL string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	return conn, nil
}

func (client *RabbitMQClient) ConsumeQueue(queueName string) (<-chan []byte, error) {
	ch, err := client.Connection.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}

	out := make(chan []byte)

	go func() {
		for d := range msgs {
			out <- d.Body
		}
		close(out)
	}()

	return out, nil
}
