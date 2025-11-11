package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func processChannel[T any](c <-chan amqp.Delivery, handler func(T)) {
	for d := range c {
		var data T
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		handler(data)
		if err = d.Ack(false); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	newChan, NewQ, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	dChan, err := newChan.Consume(NewQ.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go processChannel(dChan, handler)
	return nil
}
