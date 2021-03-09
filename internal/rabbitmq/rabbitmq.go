package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Service interface {
	Connect() error
	Publish(message string) error
	Consume()
}

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func (r *RabbitMQ) Connect() error {
	fmt.Println("Connecting to rabbitmq")
	var err error
	r.Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		return err
	}

	fmt.Println("Connected")

	r.Channel, err = r.Conn.Channel()
	if err != nil {
		return err
	}

	_, err = r.Channel.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false, nil,
	)

	return nil
}

// Consumes - consumes messages from our test queue
func (r *RabbitMQ) Consume() {

	msgs, err := r.Channel.Consume(
		"TestQueue",
		"",
		true, false, false, false, nil,
	)

	if err != nil {
		fmt.Println(err)
	}

	for msg := range msgs {
		fmt.Printf("Received message %s", msg.Body)
	}
}

// Publish - takes in a string message and publishes to a queue
func (r *RabbitMQ) Publish(message string) error {
	err := r.Channel.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	fmt.Println("Published message")
	return err
}

func NewRabbitMQService() *RabbitMQ {
	return &RabbitMQ{}
}
