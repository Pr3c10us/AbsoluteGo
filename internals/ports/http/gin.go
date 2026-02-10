package http

import (
	"fmt"

	"github.com/Pr3c10us/absolutego/internals/adapters"
	"github.com/Pr3c10us/absolutego/internals/ports/http/book"
	"github.com/Pr3c10us/absolutego/internals/ports/http/script"
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
	ginServer.scriptRoutes()

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
		bookRoute.GET("", handler.GetBooks)
		bookRoute.DELETE("/:id", handler.DeleteBook)

		bookRoute.GET("/page", handler.GetPages)
		bookRoute.GET("/panel", handler.GetPanels)

		bookRoute.POST("/chapter", handler.AddChapter)
		bookRoute.GET("/chapter", handler.GetChapters)
		bookRoute.DELETE("/chapter/:id", handler.DeleteChapter)

	}
}

func (server *GinServer) scriptRoutes() {
	handler := script.NewScriptHandler(server.Services.ScriptServices)
	scriptRoute := server.Engine.Group("/api/v1/script")
	{
		scriptRoute.GET("", handler.GetScripts)
		scriptRoute.POST("", handler.GenerateScripts)
		scriptRoute.DELETE("/:id", handler.DeleteScript)

		scriptRoute.GET("/split/:scriptId", handler.GetSplits)
		scriptRoute.POST("/split/:scriptId", handler.GenerateSplits)
		scriptRoute.DELETE("/split/:scriptId", handler.DeleteSplits)
	}
}

func (server *GinServer) Run() {
	err := server.Engine.Run(server.Environment.Port)
	if err != nil {
		fmt.Println("panic", "failed to start server")
	}
}
