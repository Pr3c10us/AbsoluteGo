package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/ports/queue/generateaudio"
	"sync"

	"github.com/Pr3c10us/absolutego/internals/adapters"
	event2 "github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/queueport"
	"github.com/Pr3c10us/absolutego/internals/ports/queue/addchapter"
	"github.com/Pr3c10us/absolutego/internals/ports/queue/generatescript"
	"github.com/Pr3c10us/absolutego/internals/ports/queue/generatesplits"
	"github.com/Pr3c10us/absolutego/internals/services"
	"github.com/Pr3c10us/absolutego/packages/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

const maxRetries = 5
const retryHeaderKey = "x-retry-count"

type AMQConsumer struct {
	environmentVariables *configs.EnvironmentVariables
	amqp                 *amqp.Channel
	Listen               chan struct{}
	services             *services.Services
	adapter              *adapters.Adapters
}

func NewAMQConsumer(environmentVariables *configs.EnvironmentVariables, amqp *amqp.Channel, services *services.Services, adapters *adapters.Adapters) *AMQConsumer {
	amqConsumer := &AMQConsumer{
		environmentVariables: environmentVariables,
		amqp:                 amqp,
		Listen:               make(chan struct{}),
		services:             services,
		adapter:              adapters,
	}

	amqConsumer.Consume(addchapter.Handler, string(queue.QueueAddChapter), 5)
	amqConsumer.Consume(generatescript.Handler, string(queue.QueueGenScript), 20)
	amqConsumer.Consume(generatesplits.Handler, string(queue.QueueGenScriptSplit), 20)
	amqConsumer.Consume(generateaudio.Handler, string(queue.QueueGenAudio), 20)

	return amqConsumer
}

func (amqConsumer *AMQConsumer) Consume(handler func(c *queueport.Context) (*queueport.HandlerResult, error), queueName string, maxWorkers int) {
	q, err := amqConsumer.amqp.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to declare queue %v: %v", queueName, err))
	}

	messages, err := amqConsumer.amqp.Consume(
		q.Name,
		"",    // consumer
		false, // auto-ack — disabled so we can manually ack/nack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to register %v consumer: %v", queueName, err))
	}

	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for message := range messages {
				retryCount := getRetryCount(message)

				var m queue.Message
				err = json.Unmarshal(message.Body, &m)
				if err != nil {
					_ = message.Ack(false)
					continue
				}

				eventId := m.EventId
				var event *event2.Event
				event, err = amqConsumer.adapter.EventImplementation.GetEvent(eventId)
				if err != nil || event == nil {
					_ = message.Ack(false)
					continue
				}

				_ = amqConsumer.adapter.EventImplementation.Update(eventId, event2.Event{
					Status: event2.StatusProcessing,
				})

				ctx := queueport.Context{
					Adapters:             amqConsumer.adapter,
					Services:             amqConsumer.services,
					Data:                 m.Data,
					EnvironmentVariables: amqConsumer.environmentVariables,
				}

				var result *queueport.HandlerResult
				if result, err = handler(&ctx); err != nil {
					fmt.Println(err)
					var chapterId, scriptId, vabId int64
					if result != nil {
						chapterId = result.ChapterId
						scriptId = result.ScriptId
						vabId = result.VabId
					}

					if retryCount >= maxRetries {
						_ = amqConsumer.adapter.EventImplementation.Update(eventId, event2.Event{
							Status:    event2.StatusFailed,
							ChapterId: chapterId,
							ScriptId:  scriptId,
							VabId:     vabId,
						})
						_ = message.Ack(false)
						continue
					}

					if pubErr := amqConsumer.republish(queueName, message, retryCount+1); pubErr != nil {
						_ = amqConsumer.adapter.EventImplementation.Update(eventId, event2.Event{
							Status:    event2.StatusFailed,
							ChapterId: chapterId,
							ScriptId:  scriptId,
							VabId:     vabId,
						})
						_ = message.Nack(false, false)
						continue
					}

					_ = amqConsumer.adapter.EventImplementation.Update(eventId, event2.Event{
						Status:    event2.StatusRetry,
						ChapterId: chapterId,
						ScriptId:  scriptId,
						VabId:     vabId,
					})
					_ = message.Ack(false)
				} else {
					_ = amqConsumer.adapter.EventImplementation.Update(eventId, event2.Event{
						Status:    event2.StatusSuccessful,
						ChapterId: result.ChapterId,
						ScriptId:  result.ScriptId,
						VabId:     result.VabId,
					})
					_ = message.Ack(false)
				}
			}
		}()
	}
}

func getRetryCount(msg amqp.Delivery) int {
	if msg.Headers == nil {
		return 0
	}

	val, ok := msg.Headers[retryHeaderKey]
	if !ok {
		return 0
	}

	switch v := val.(type) {
	case int32:
		return int(v)
	case int64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}

func (amqConsumer *AMQConsumer) republish(queueName string, original amqp.Delivery, retryCount int) error {
	headers := original.Headers
	if headers == nil {
		headers = amqp.Table{}
	}
	headers[retryHeaderKey] = int32(retryCount)

	return amqConsumer.amqp.PublishWithContext(
		context.Background(),
		"",        // default exchange
		queueName, // routing key = queue name
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			Headers:     headers,
			ContentType: original.ContentType,
			Body:        original.Body,
		},
	)
}
