package pubsub

import (
	"context"
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	TransientQueue SimpleQueueType = "transient"
	DurableQueue   SimpleQueueType = "durable"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx := context.Background()
	ch.PublishWithContext(ctx, exchange, key, false, false, amqp.Publishing{
		ContentType: "applcaiton/json",
		Body:        data,
	})
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	newChan, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	qDurable := false
	qAutoDelete := false
	qExclusive := false

	switch queueType {
	case TransientQueue:
		qAutoDelete = true
		qExclusive = true
	case DurableQueue:
		qDurable = true
	default:
		return nil, amqp.Queue{}, errors.New("unsupported queue type")
	}

	newQueue, err := newChan.QueueDeclare(queueName, qDurable, qAutoDelete, qExclusive, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = newChan.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return newChan, newQueue, nil
}
