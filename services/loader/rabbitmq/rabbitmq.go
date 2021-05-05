package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vladimir3322/stonent_go/config"
	"log"
)

//var publishConnection, _ = amqp.Dial(config.RabbitMQUrl)
//var publishAmqpChannel, _ = publishConnection.Channel()
//var publishQueue, _ = publishAmqpChannel.QueueDeclare(
//	config.RabbitMQQueueName,
//	true,
//	false,
//	false,
//	false,
//	nil,
//)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func PushEvent(data []byte) {
	conn, err := amqp.Dial(config.RabbitMQUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		config.RabbitMQQueueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	erro := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})

	if erro != nil {
		fmt.Println(err)
	}
}

func ConsumeEvents() {
	conn, err := amqp.Dial(config.RabbitMQUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		config.RabbitMQQueueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}
