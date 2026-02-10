package script

import (
	script2 "github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/services/script"
	"github.com/Pr3c10us/absolutego/internals/services/script/commands"
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
		BookID int64   `form:"bookId" binding:"omitempty,gt=0"`
		Name   string  `form:"name"   binding:"omitempty,min=1,max=255"`
		Ids    []int64 `form:"id"   binding:"omitempty,dive,gt=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	var scripts []script2.Script
	var err error
	if scripts, err = h.service.GetScripts.Handle(req.BookID, req.Name, req.Ids); err != nil {
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

func (h *Handler) GenerateScripts(c *gin.Context) {
	var req struct {
		BookId          int64   `json:"bookId" binding:"required,gt=0"`
		Name            string  `json:"name"   binding:"required,min=1,max=255"`
		Chapters        []int   `json:"chapters" binding:"required,min=1"`
		PreviousScripts []int64 `json:"previousScripts"   binding:"omitempty,dive,gt=0"`
	}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	var scriptContent string
	var scriptId int64
	var err error
	if scriptContent, scriptId, err = h.service.GenerateScript.Handle(commands.Parameters{
		BookId:          req.BookId,
		Name:            req.Name,
		Chapters:        req.Chapters,
		PreviousScripts: req.PreviousScripts,
	}); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"script": scriptContent, "scriptId": scriptId}, nil).Send(c)
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
