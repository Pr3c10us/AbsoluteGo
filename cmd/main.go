package main

import (
	"github.com/Pr3c10us/absolutego/internals/adapters"
	"github.com/Pr3c10us/absolutego/internals/ports"
	"github.com/Pr3c10us/absolutego/internals/services"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
)

var (
	environmentVariables = configs.LoadEnvironment()
)

func main() {
	newS3Client := utils.NewS3Client(environmentVariables)
	newSQLClient := utils.NewSQLClient(environmentVariables)

	adapterDependencies := adapters.AdapterDependencies{
		EnvironmentVariables: environmentVariables,
		S3Client:             newS3Client,
		DB:                   newSQLClient,
	}
	newAdapters := adapters.NewAdapters(adapterDependencies)
	newServices := services.NewServices(newAdapters)
	newPort := ports.NewPorts(newAdapters, newServices, environmentVariables)
	newPort.GinServer.Run()
}
