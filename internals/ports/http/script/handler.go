package script

import (
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
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
		Page   int     `form:"page" binding:"omitempty,min=1"`
		Limit  int     `form:"limit" binding:"omitempty,min=1,max=100"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	var scripts []script2.Script
	var err error
	if scripts, err = h.service.GetScripts.Handle(req.BookID, req.Name, req.Ids, req.Page, req.Limit); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"scripts": scripts}, nil).Send(c)
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

	var err error
	if err = h.service.CreateScript.Handle(commands.CreateScriptParameters{
		BookId:          req.BookId,
		Name:            req.Name,
		Chapters:        req.Chapters,
		PreviousScripts: req.PreviousScripts,
	}); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("added to queue", nil, nil).Send(c)
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

func (h *Handler) GetSplits(c *gin.Context) {
	var req struct {
		ScriptID int64 `uri:"scriptId" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
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

func (h *Handler) GenerateSplits(c *gin.Context) {
	var req struct {
		ScriptId int64 `uri:"scriptId" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.CreateSplits.Handle(req.ScriptId); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"message": "added to queue"}, nil).Send(c)
}

func (h *Handler) GenerateAudios(c *gin.Context) {
	var body struct {
		ScriptId   int64  `json:"scriptId" binding:"required,gt=0"`
		Voice      string `json:"voice" binding:"required,oneof=Zephyr Puck Charon Kore Fenrir Leda Orus Aoede Callirrhoe Autonoe Enceladus Iapetus Umbriel Algieba Despina Erinome Algenib Rasalgethi Laomedeia Achernar Alnilam Schedar Gacrux Pulcherrima Achird Zubenelgenubi Vindemiatrix Sadachbia Sadaltager Sulafat"`
		VoiceStyle string `json:"voiceStyle" binding:"omitempty"`
	}
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.CreateAudios.Handle(body.ScriptId, ai.Voice(body.Voice), body.VoiceStyle); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"message": "added to queue"}, nil).Send(c)
}

func (h *Handler) GenerateAudio(c *gin.Context) {
	var body struct {
		SplitId    int64  `json:"splitId" binding:"required,gt=0"`
		Voice      string `json:"voice" binding:"required,oneof=Zephyr Puck Charon Kore Fenrir Leda Orus Aoede Callirrhoe Autonoe Enceladus Iapetus Umbriel Algieba Despina Erinome Algenib Rasalgethi Laomedeia Achernar Alnilam Schedar Gacrux Pulcherrima Achird Zubenelgenubi Vindemiatrix Sadachbia Sadaltager Sulafat"`
		VoiceStyle string `json:"voiceStyle" binding:"omitempty"`
	}
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.CreateAudio.Handle(body.SplitId, ai.Voice(body.Voice), body.VoiceStyle); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"message": "added to queue"}, nil).Send(c)
}

func (h *Handler) DeleteSplits(c *gin.Context) {
	var req struct {
		ScriptId int64 `uri:"scriptId" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.DeleteSplits.Handle(req.ScriptId); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("splits deleted", nil, nil).Send(c)
}

func (h *Handler) GenerateVideos(c *gin.Context) {
	var body struct {
		ScriptId int64 `uri:"scriptId" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&body); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.CreateVideos.Handle(body.ScriptId); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"message": "added to queue"}, nil).Send(c)
}

func (h *Handler) GenerateVideo(c *gin.Context) {
	var body struct {
		SplitId int64 `uri:"splitId" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&body); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.CreateVideo.Handle(body.SplitId); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"message": "added to queue"}, nil).Send(c)
}
