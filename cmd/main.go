package main

import (
	"github.com/Pr3c10us/absolutego/internals/adapters"
	"github.com/Pr3c10us/absolutego/internals/ports"
	"github.com/Pr3c10us/absolutego/internals/services"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	environmentVariables = configs.LoadEnvironment()
)

func main() {
	newS3Client := utils.NewS3Client(environmentVariables)
	newSQLClient := utils.NewSQLClient(environmentVariables)
	newGoogleGenAIClient := utils.NewGoogleGenAIClient(environmentVariables)

	newAMQConnection := utils.NewAMQConnection(environmentVariables)
	defer func(amqp *amqp.Connection) {
		_ = amqp.Close()
	}(newAMQConnection)
	newAMQChannel := utils.NewAMQChannel(newAMQConnection)
	defer func(newAMQChannel *amqp.Channel) {
		_ = newAMQChannel.Close()
	}(newAMQChannel)

	adapterDependencies := adapters.AdapterDependencies{
		EnvironmentVariables: environmentVariables,
		GoogleGenAIClient:    newGoogleGenAIClient,
		S3Client:             newS3Client,
		DB:                   newSQLClient,
		AMQP:                 newAMQChannel,
	}
	newAdapters := adapters.NewAdapters(adapterDependencies)
	newServices := services.NewServices(newAdapters)
	newPort := ports.NewPorts(newAdapters, newServices, environmentVariables)
	newPort.GinServer.Run()
}
