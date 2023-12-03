package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

var (
	rabbit_host     = os.Getenv("RABBIT_HOST")
	rabbit_port     = os.Getenv("RABBIT_PORT")
	rabbit_username = os.Getenv("RABBIT_USERNAME")
	rabbit_password = os.Getenv("RABBIT_PASSWORD")
	api_port        = os.Getenv("PORT")
)

func main() {
	log.Println("===================== ENV Variables =======================")
	log.Printf("rabbit host ➜ %v", rabbit_host)
	log.Printf("rabbit port ➜ %v", rabbit_port)
	log.Printf("rabbit username	➜ %v", rabbit_username)
	log.Printf("rabbit password	➜ %v", rabbit_password)
	log.Printf("api port ➜ %v", api_port)
	log.Println("===========================================================")

	h := gin.Default()
	srv := &http.Server{
		Handler: h,
		Addr:    fmt.Sprintf("0.0.0.0:%v", api_port),
	}

	conn, err := NewRabbitMqConnection(rabbit_username, rabbit_password, rabbit_host, rabbit_port)
	HandleMainErrors(err, conn, nil)

	ch, err := NewRabbitMqChannel(conn)
	HandleMainErrors(err, conn, ch)

	queue, err := NewQueue(ch, RabbitQueue{
		Name:             "publisher-queue",
		Durable:          false,
		DeleteWhenUnused: false,
		Exclusive:        false,
		NoWait:           false,
	})
	HandleMainErrors(err, conn, ch)

	h.POST("/produce/:msg", func(c *gin.Context) {
		msg := c.Param("msg")
		err := submit(*queue, *ch, msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"response": "published",
		})
	})

	h.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"response": "up",
		})
	})

	run(srv)
}

func HandleMainErrors(err error, conn *amqp.Connection, ch *amqp.Channel) {
	if err != nil {
		log.Println("shutting down the server .. ")
		ch.Close()
		conn.Close()
		os.Exit(1)
	}
}

func run(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error trying to start the server ➜ %v", err)
	}
}

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

func submit(queue amqp.Queue, ch amqp.Channel, msg string) error {
	err := ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		log.Printf("error trying to publish a msg ➜ %v", err)
		return err
	}
	return nil
}


