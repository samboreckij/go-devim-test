package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	msg "./my_message"
	"github.com/go-snorlax/env"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"integer_div", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")
	/*
		body := "Hello World!"
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body)
			})
		log.Printf(" [x] Sent %s", body)
		failOnError(err, "Failed to publish a message") */
	//Sending messages
	forever := make(chan bool)

	go func() {

		for d := range msgs {
			//log.Printf("Received a message: %s", d.Body)
			numMessage := &msg.Message{}
			if err := proto.Unmarshal(d.Body, numMessage); err != nil {
				log.Fatalln("Failed to parse message:", err)
			}
			if numMessage.GetNumber()%denominator == 0 {
				addToStorage(counter, numMessage.GetNumber(), storage)
			}
		}
		viewStorage(storage)
	}()

	//log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
