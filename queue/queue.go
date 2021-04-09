package queue

import (
	"sync"

	"github.com/streadway/amqp"
)

var onceQueue, onceChannel sync.Once
var conn *amqp.Connection = nil
var channel *amqp.Channel = nil

func queueConnect() *amqp.Connection {
	onceQueue.Do(func() {
		var err error
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			panic(err)
		}
	})
	return conn
}

func channelConnect() *amqp.Channel {
	onceChannel.Do(func() {
		conn := queueConnect()
		var err error
		channel, err = conn.Channel()
		if err != nil {
			panic(err)
		}
	})

	return channel
}

func getSubmissionQueue() amqp.Queue {
	ch := channelConnect()
	queue, err := ch.QueueDeclare(
		"submissions", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		panic(err)
	}
	return queue
}

func SendMessage(body string) error {
	ch := channelConnect()
	q := getSubmissionQueue()
	err := ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	return err
}
