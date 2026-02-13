package queue

import (
	"context"
	"fmt"
	"github.com/Pr3c10us/fanatix/internals/domains/queue"
	"github.com/Pr3c10us/fanatix/packages/configs"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type AMQRepository struct {
	channel             *amqp.Channel
	environmentVariable *configs.EnvironmentVariables
}

func NewAMQRepository(channel *amqp.Channel, environmentVariable *configs.EnvironmentVariables) queue.Repository {
	return &AMQRepository{
		channel,
		environmentVariable,
	}
}

func (repo *AMQRepository) Publish(params *queue.MessageParams) error {
	q, err := repo.channel.QueueDeclare(
		params.Queue, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = repo.channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(params.Message),
		})
	if err != nil {
		return err
	}
	return nil
}

type KafkaRepository struct {
	producer            *kafka.Producer
	environmentVariable *configs.EnvironmentVariables
}

func NewKafkaRepository(producer *kafka.Producer, environmentVariable *configs.EnvironmentVariables) queue.Repository {
	return &KafkaRepository{
		producer,
		environmentVariable,
	}
}

func (repo *KafkaRepository) Publish(params *queue.MessageParams) error {
	err := repo.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic: &params.Queue, Partition: kafka.PartitionAny,
			},
			Value: []byte(params.Message),
			Key:   []byte(params.Key),
		},

		nil, // delivery channel
	)
	if err != nil {
		return err
	}

	go func() {
		for e := range repo.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					//fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
					//	*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	return nil
}
