package queue

import (
	"sync"

	"github.com/streadway/amqp"
)

var onceQueue, onceChannel sync.Once
var queueError, channelError error
var conn *amqp.Connection = nil
var channel *amqp.Channel = nil

func queueConnect() (*amqp.Connection, error) {
	onceQueue.Do(func() {
		conn, queueError = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if queueError != nil {
			return
		}
	})
	if queueError != nil {
		onceQueue = sync.Once{}
	}
	return conn, queueError
}

func channelConnect() (*amqp.Channel, error) {
	onceChannel.Do(func() {
		conn, channelError = queueConnect()
		if channelError != nil {
			return
		}
		channel, channelError = conn.Channel()
		if channelError != nil {
			return
		}
	})
	if channelError != nil {
		onceChannel = sync.Once{}
	}
	return channel, channelError
}

func getSubmissionQueue() (amqp.Queue, error) {
	ch, err := channelConnect()
	if err != nil {
		return amqp.Queue{}, err
	}
	queue, err := ch.QueueDeclare(
		"submissions", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	return queue, err
}

func Submit(id string) error {
	ch, err := channelConnect()
	if err != nil {
		return err
	}
	q, err := getSubmissionQueue()
	if err != nil {
		return err
	}
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(id),
		})
	return err
}

func GetSubmissionChannel() (<-chan amqp.Delivery, error) {
	ch, err := channelConnect()
	if err != nil {
		return nil, err
	}
	getSubmissionQueue()
	msgs, err := ch.Consume(
		"submissions",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return msgs, err
}
