package ports

import (
	"github.com/Pr3c10us/absolutego/internals/adapters"
	"github.com/Pr3c10us/absolutego/internals/ports/http"
	"github.com/Pr3c10us/absolutego/internals/services"
	"github.com/Pr3c10us/absolutego/packages/configs"
)

type Ports struct {
	GinServer *http.GinServer
}

func NewPorts(adapters *adapters.Adapters, services *services.Services, environment *configs.EnvironmentVariables) *Ports {
	return &Ports{
		GinServer: http.NewGinServer(services, adapters, environment),
	}
}
