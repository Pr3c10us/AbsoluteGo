package queueport

import (
	"github.com/Pr3c10us/fanatix/internals/infrastructures/adapters"
	"github.com/Pr3c10us/fanatix/internals/services"
	"github.com/Pr3c10us/fanatix/packages/configs"
	"github.com/Pr3c10us/fanatix/packages/logger"
	"github.com/rabbitmq/amqp091-go"
)

type Context struct {
	Adapters             *adapters.Adapters
	Services             *services.Services
	Logger               logger.Logger
	AMQMessage           *amqp091.Delivery
	BytesMessage         []byte
	EnvironmentVariables *configs.EnvironmentVariables
}
