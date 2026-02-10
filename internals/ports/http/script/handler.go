package script

import (
	script2 "github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/services/script"
	"github.com/Pr3c10us/absolutego/packages/response"
	"github.com/Pr3c10us/absolutego/packages/validator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service script.Services
}

func NewScriptHandler(service script.Services) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) GetScripts(c *gin.Context) {
	var req struct {
		BookID int64  `form:"bookId" binding:"omitempty,gt=0"`
		Name   string `form:"name"   binding:"omitempty,min=1,max=255"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	var scripts []script2.Script
	var err error
	if scripts, err = h.service.GetScripts.Handle(req.BookID, req.Name); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"scripts": scripts}, nil).Send(c)
}

func (h *Handler) GetSplits(c *gin.Context) {
	var req struct {
		ScriptID int64 `form:"scriptId" binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	var splits []script2.Split
	var err error
	if splits, err = h.service.GetSplits.Handle(req.ScriptID); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"splits": splits}, nil).Send(c)
}

func (h *Handler) DeleteScript(c *gin.Context) {
	var uri struct {
		Id int64 `uri:"id" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.DeleteScript.Handle(uri.Id); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("script deleted", nil, nil).Send(c)
}
