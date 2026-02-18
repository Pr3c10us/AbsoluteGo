package vab

import (
	vab2 "github.com/Pr3c10us/absolutego/internals/domains/vab"

	"github.com/Pr3c10us/absolutego/internals/services/vab"
	"github.com/Pr3c10us/absolutego/packages/response"
	"github.com/Pr3c10us/absolutego/packages/validator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service vab.Services
}

func NewVABHandler(service vab.Services) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) CreateVAB(c *gin.Context) {
	var req struct {
		Name     string `form:"name" binding:"required,min=1,max=255"`
		ScriptId int64  `form:"scriptId" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.QueueVAB.Handle(req.ScriptId, req.Name); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("vab queued", nil, nil).Send(c)
}

func (h *Handler) GetVABs(c *gin.Context) {
	var req struct {
		Name     string `form:"name" binding:"omitempty,min=1,max=255"`
		ScriptId int64  `form:"scriptId" binding:"omitempty,min=1"`
		BookId   int64  `form:"bookId" binding:"omitempty,min=1"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}
	var vabs []vab2.VAB
	var err error
	if vabs, err = h.service.GetVABs.Handle(req.ScriptId, req.BookId, req.Name); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"vabs": vabs}, nil).Send(c)
}

func (h *Handler) DeleteVAB(c *gin.Context) {
	var uri struct {
		Id int64 `uri:"id" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.DeleteVAB.Handle(uri.Id, 0); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("vab deleted", nil, nil).Send(c)
}
