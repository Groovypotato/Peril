package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
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
