package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	go server()
	go client()

	var a string
	fmt.Scanln(&a)
}

func getQeue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "Failed to conect to RabbitMq")

	ch, err := conn.Channel()
	failOnError(err, "Fialed to open channel")

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnError(err, "Failed to declare queue")

	return conn, ch, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func server() {
	conn, ch, q := getQeue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello Bunny"),
	}

	for {
		ch.Publish("", q.Name, false, false, msg)
	}

}

func client() {
	_, ch, q := getQeue()
	defer ch.Close()

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		log.Printf("Received message with message: %s", msg.Body)
	}
}
