package queue

import (
	"context"
	"encoding/json"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/packages/configs"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type AMQImplementation struct {
	channel             *amqp.Channel
	environmentVariable *configs.EnvironmentVariables
}

func NewAMQImplementation(channel *amqp.Channel, environmentVariable *configs.EnvironmentVariables) queue.Interface {
	return &AMQImplementation{
		channel,
		environmentVariable,
	}
}

func (repo *AMQImplementation) Publish(params *queue.MessageParams) error {
	q, err := repo.channel.QueueDeclare(
		string(params.Queue), // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var messageByte []byte
	messageByte, err = json.Marshal(params.Message)
	if err != nil {
		return err
	}

	err = repo.channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageByte,
		})
	if err != nil {
		return err
	}
	return nil
}
