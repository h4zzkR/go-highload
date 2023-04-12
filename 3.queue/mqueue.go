package main

import (
	"context"
	"log"
	"time"

	email "queue/email"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	conn amqp.Connection
	ch   *amqp.Channel
}

func newMessageQueue() *Queue {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"notification", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)

	failOnError(err, "Failed to declare an exchange")

	return &Queue{
		conn: *conn,
		ch:   ch,
	}
}

func (q *Queue) sendEmail(item email.MailItem) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	encoded, err := MailItem2JSON(item)
	failOnError(err, "Failed to convert message to JSON format")

	err = q.ch.PublishWithContext(ctx,
		"notification", // exchange
		"email",        // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        encoded,
		})

	failOnError(err, "Failed to publish a message")
	log.Printf("Sent %s\n", encoded)
}
