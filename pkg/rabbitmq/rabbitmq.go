package rabbitmq

import (
	"errors"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var (
	rabbit_host     = os.Getenv("RABBIT_HOST")
	rabbit_port     = os.Getenv("RABBIT_PORT")
	rabbit_username = os.Getenv("RABBIT_USERNAME")
	rabbit_password = os.Getenv("RABBIT_PASSWORD")
	api_port        = os.Getenv("PORT")
)

func NewRabbitMqConnection(username, password, host, port string) (*amqp.Connection, error) {
	connString := "amqp://" + rabbit_username + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Printf("error trying to connect to rabbitmq ➜ %v", err)
		return nil, err
	}
	if conn.IsClosed() {
		return nil, errors.New("connection is closed")
	}
	return conn, nil
}

func NewRabbitMqChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("error trying to open a channel to the queue ➜ %v", err)
		return nil, err
	}
	return ch, nil
}

type RabbitQueue struct {
	Name             string
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	NoWait           bool
}

func NewQueue(ch *amqp.Channel, queueSpecs RabbitQueue) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueSpecs.Name,
		queueSpecs.Durable,
		queueSpecs.DeleteWhenUnused,
		queueSpecs.Exclusive,
		queueSpecs.NoWait,
		nil,
	)
	if err != nil {
		log.Printf("error trying to create a queue ➜ %v", err)
		return nil, err
	}
	return &q, nil
}
