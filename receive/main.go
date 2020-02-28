package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	"./env"
	msg "./my_message"
	"github.com/streadway/amqp"
)

func viewStorage(storage map[int64]int64) {
	fmt.Println(storage)
}

func addToStorage(count int64, number int64, storage map[int64]int64) {
	storage[count] = number
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	//Variables
	var denominator int64 = int64(env.Get("DENOMINATOR_VALUE").Int(7))
	var counter int64 = 1
	var storage = make(map[int64]int64)

	//Trying connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare( //create queue
		"integer_div", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
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

	//Receiving messages
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

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
