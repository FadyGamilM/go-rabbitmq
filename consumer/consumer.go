package main

import (
	"log"
	"os"

	"github.com/FadyGamilM/gorabbitmq/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

var (
	rabbit_host     = os.Getenv("RABBIT_HOST")
	rabbit_port     = os.Getenv("RABBIT_PORT")
	rabbit_username = os.Getenv("RABBIT_USERNAME")
	rabbit_password = os.Getenv("RABBIT_PASSWORD")
)

func main() {
	log.Println("===================== ENV Variables =======================")
	log.Printf("rabbit host ➜ %v", rabbit_host)
	log.Printf("rabbit port ➜ %v", rabbit_port)
	log.Printf("rabbit username	➜ %v", rabbit_username)
	log.Printf("rabbit password	➜ %v", rabbit_password)
	log.Println("===========================================================")

	conn, err := rabbitmq.NewRabbitMqConnection(rabbit_username, rabbit_password, rabbit_host, rabbit_port)
	HandleMainErrors(err, conn, nil)

	ch, err := rabbitmq.NewRabbitMqChannel(conn)
	HandleMainErrors(err, conn, ch)

	queue, err := rabbitmq.NewQueue(ch, rabbitmq.RabbitQueue{
		Name:             "publisher-queue",
		Durable:          false,
		DeleteWhenUnused: false,
		Exclusive:        false,
		NoWait:           false,
	})
	HandleMainErrors(err, conn, ch)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false, // -> disable auto ack so i will ack manually
		false,
		false,
		false,
		nil,
	)
	HandleMainErrors(err, conn, ch)

	// spin a go routine running forever listening for msgs and consuming them
	consumerRoutine := make(chan bool)

	go func() {
		for msg := range msgs {
			log.Printf("consumed msg ➜ %v\n", string(msg.Body))
			// manually ack
			msg.Ack(false) // not true so i will avoid acknowledgina ll prior unacknowledged msgs
		}
	}()

	log.Println("consumer is running on background .. ")
	<-consumerRoutine
}

func HandleMainErrors(err error, conn *amqp.Connection, ch *amqp.Channel) {
	if err != nil {
		log.Println("shutting down the server .. ")
		ch.Close()
		conn.Close()
		os.Exit(1)
	}
}
