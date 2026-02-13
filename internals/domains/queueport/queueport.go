package queueport

import (
	"github.com/Pr3c10us/absolutego/internals/adapters"
	"github.com/Pr3c10us/absolutego/internals/services"
	"github.com/Pr3c10us/absolutego/packages/configs"
)

type Context struct {
	Adapters             *adapters.Adapters
	Services             *services.Services
	Data                 []byte
	BytesMessage         []byte
	EnvironmentVariables *configs.EnvironmentVariables
}
