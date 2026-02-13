package queue

import (
	"fmt"
	"github.com/Pr3c10us/fanatix/internals/domains/queueport"
	"github.com/Pr3c10us/fanatix/internals/infrastructures/adapters"
	"github.com/Pr3c10us/fanatix/internals/infrastructures/ports/queue/events"
	"github.com/Pr3c10us/fanatix/internals/infrastructures/ports/queue/livescores"
	"github.com/Pr3c10us/fanatix/internals/services"
	"github.com/Pr3c10us/fanatix/packages/configs"
	"github.com/Pr3c10us/fanatix/packages/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQConsumer struct {
	logger               logger.Logger
	environmentVariables *configs.EnvironmentVariables
	amqp                 *amqp.Channel
	Listen               chan struct{}
	services             *services.Services
	adapter              *adapters.Adapters
}

func NewAMQConsumer(newLogger logger.Logger, environmentVariables *configs.EnvironmentVariables, amqp *amqp.Channel, services *services.Services, adapters *adapters.Adapters) *AMQConsumer {
	amqConsumer := &AMQConsumer{
		logger:               newLogger,
		environmentVariables: environmentVariables,
		amqp:                 amqp,
		Listen:               make(chan struct{}),
		services:             services,
		adapter:              adapters,
	}

	// live score
	amqConsumer.Consume(livescores.Handler, environmentVariables.LiveScoresQueue)
	amqConsumer.Consume(events.Handler, environmentVariables.KafkaTopics.EventsTopic)

	return amqConsumer
}

func (amqConsumer *AMQConsumer) Consume(handler func(c *queueport.Context), queueName string) {
	q, err := amqConsumer.amqp.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		amqConsumer.logger.LogWithFields("panic", fmt.Sprintf("Failed to declare queue %v", queueName), err)
	}

	messages, err := amqConsumer.amqp.Consume(
		q.Name,
		"",    // amqConsumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		amqConsumer.logger.LogWithFields("panic", fmt.Sprintf("Failed to register %v amqConsumer", queueName), err)
	}

	go func() {
		for message := range messages {
			context := queueport.Context{
				Adapters:             amqConsumer.adapter,
				Services:             amqConsumer.services,
				Logger:               amqConsumer.logger,
				AMQMessage:           &message,
				EnvironmentVariables: amqConsumer.environmentVariables,
			}
			handler(&context)
		}
	}()

	amqConsumer.logger.Log("debug", fmt.Sprintf("[*] Waiting for %v messages.", queueName))
}
