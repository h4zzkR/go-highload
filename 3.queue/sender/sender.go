package main

import (
	"encoding/json"
	"fmt"
	"log"

	email "queue/email"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func prepareDeliveryService(ch *amqp.Channel) <-chan amqp.Delivery {

	q, err := ch.QueueDeclare(
		"emails", // name
		false,    // durable
		false,    // delete when unused
		true,     // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue for mail")

	err = ch.QueueBind(
		q.Name,         // queue name
		"email",        // routing key
		"notification", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a mail queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a mail consumer")

	return msgs
}

func delivery(dService <-chan amqp.Delivery) {
	for m := range dService {
		mail := email.MailItem{}
		err := json.Unmarshal(m.Body, &mail)
		failOnError(err, "Failed to Read Message")
		fmt.Print(mail)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

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

	dService := prepareDeliveryService(ch)

	forever := make(chan bool)

	go delivery(dService)

	fmt.Println("Mail delivery service is running...")

	<-forever
}
