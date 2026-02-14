package event

import (
	event2 "github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/services/event"
	"github.com/Pr3c10us/absolutego/packages/response"
	"github.com/Pr3c10us/absolutego/packages/validator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services event.Services
}

func NewEventHandler(service event.Services) Handler {
	return Handler{
		services: service,
	}
}

func (handler *Handler) GetEvents(context *gin.Context) {
	var filter event2.Filter
	if err := context.ShouldBindQuery(&filter); err != nil {
		_ = context.Error(validator.ValidateRequest(err))
		return
	}

	events, err := handler.services.GetEvents.Handle(filter)
	if err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"events": events}, nil).Send(context)
}
