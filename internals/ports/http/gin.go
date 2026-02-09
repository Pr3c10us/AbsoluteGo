package http

import (
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/adapters"
	"github.com/Pr3c10us/absolutego/internals/ports/http/book"
	"github.com/Pr3c10us/absolutego/internals/services"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/middlewares"
	"github.com/Pr3c10us/absolutego/packages/response"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	Services    *services.Services
	Adapters    *adapters.Adapters
	Engine      *gin.Engine
	Environment *configs.EnvironmentVariables
}

func NewGinServer(services *services.Services, adapters *adapters.Adapters, environment *configs.EnvironmentVariables) *GinServer {
	ginServer := &GinServer{
		Services:    services,
		Adapters:    adapters,
		Engine:      gin.Default(),
		Environment: environment,
	}

	// Set up CORS
	ginServer.Engine.Use(middlewares.CORSMiddleware(environment.AllowedOrigins))

	// Middlewares
	ginServer.Engine.Use(middlewares.ErrorHandlerMiddleware())
	ginServer.Engine.NoRoute(middlewares.RouteNotFoundMiddleware())

	ginServer.health()
	ginServer.bookRoutes()

	return ginServer
}

func (server *GinServer) health() {
	server.Engine.GET("/health", func(c *gin.Context) {
		response.NewSuccessResponse("server up!!!", nil, nil).Send(c)
	})
}

func (server *GinServer) bookRoutes() {
	handler := book.NewBookHandler(server.Services.BookServices, server.Environment)
	bookRoute := server.Engine.Group("/api/v1/book")
	{
		bookRoute.POST("", handler.AddBook)
	}
}

func (server *GinServer) Run() {
	err := server.Engine.Run(server.Environment.Port)
	if err != nil {
		fmt.Println("panic", "failed to start server")
	}
}
