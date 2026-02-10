package book

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	book2 "github.com/Pr3c10us/absolutego/internals/domains/book"

	"github.com/Pr3c10us/absolutego/internals/services/book"
	"github.com/Pr3c10us/absolutego/internals/services/book/commands"
	"github.com/Pr3c10us/absolutego/packages/appError"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/response"
	"github.com/Pr3c10us/absolutego/packages/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service              book.Services
	environmentVariables *configs.EnvironmentVariables
}

func NewBookHandler(service book.Services, environmentVariables *configs.EnvironmentVariables) Handler {
	return Handler{
		service:              service,
		environmentVariables: environmentVariables,
	}
}

func (h *Handler) AddBook(c *gin.Context) {
	var req struct {
		Title string `form:"title" binding:"required,min=1,max=255"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.AddBook.Handle(req.Title); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("book created", nil, nil).Send(c)
}

func (h *Handler) GetBooks(c *gin.Context) {
	var req struct {
		Title string `form:"title" binding:"omitempty,min=1,max=255"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}
	var books []book2.Book
	var err error
	if books, err = h.service.GetBooks.Handle(req.Title); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"books": books}, nil).Send(c)
}

func (h *Handler) DeleteBook(c *gin.Context) {
	var uri struct {
		Id int64 `uri:"id" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.DeleteBook.Handle(uri.Id); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("book deleted", nil, nil).Send(c)
}

func (h *Handler) AddChapter(c *gin.Context) {
	var req struct {
		Book    *multipart.FileHeader `form:"book"    binding:"required"`
		Chapter int                   `form:"chapter" binding:"required,gt=0"`
		BookID  int64                 `form:"bookId"  binding:"required,gt=0"`
	}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := validateUploadedFile(req.Book); err != nil {
		_ = c.Error(err)
		return
	}

	rootPath := configs.GetRootPath()
	dest := filepath.Join(rootPath, "uploads", generateFileName(req.Book.Filename))

	if err := c.SaveUploadedFile(req.Book, dest); err != nil {
		_ = c.Error(err)
		return
	}

	err := h.service.Commands.AddChapter.Handle(commands.Parameter{
		File:    dest,
		Chapter: req.Chapter,
		BookId:  req.BookID,
	})
	if err != nil {
		fmt.Println(err)
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("chapter added", nil, nil).Send(c)
}

func (h *Handler) GetChapters(c *gin.Context) {
	var req struct {
		Numbers []int `form:"number" binding:"omitempty,dive,gt=0"`
		BookID  int64 `form:"bookId"  binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	var chapters []book2.Chapter
	var err error
	if chapters, err = h.service.GetChapters.Handle(req.BookID, req.Numbers); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"chapters": chapters}, nil).Send(c)
}

func (h *Handler) DeleteChapter(c *gin.Context) {
	var uri struct {
		Id int64 `uri:"id" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	if err := h.service.DeleteChapter.Handle(uri.Id); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("chapter deleted", nil, nil).Send(c)
}

func (h *Handler) GetPages(c *gin.Context) {
	var req struct {
		ChapterIds []int64 `form:"chapterId"  binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	var pages []book2.Page
	var err error
	if pages, err = h.service.GetPages.Handle(req.ChapterIds); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"pages": pages}, nil).Send(c)
}

func (h *Handler) GetPanels(c *gin.Context) {
	var req struct {
		PageId int64 `form:"pageId"  binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(validator.ValidateRequest(err))
		return
	}

	var panels []book2.Panel
	var err error
	if panels, err = h.service.GetPanels.Handle(req.PageId); err != nil {
		_ = c.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"panels": panels}, nil).Send(c)
}

var allowedExtensions = map[string]struct{}{
	".pdf": {},
	".cbr": {},
	".cbz": {},
	".cb7": {},
}

const maxFileSize = 5000 << 20

func validateUploadedFile(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedExtensions[ext]; !ok {
		return appError.BadRequest(fmt.Errorf(
			"unsupported file type '%s'; allowed: %s",
			ext, allowedExtensionsList(),
		))
	}
	if file.Size > maxFileSize {
		return appError.BadRequest(fmt.Errorf(
			"file size %d bytes exceeds the %d byte limit",
			file.Size, maxFileSize,
		))
	}
	return nil
}

func allowedExtensionsList() string {
	exts := make([]string, 0, len(allowedExtensions))
	for k := range allowedExtensions {
		exts = append(exts, k)
	}
	return strings.Join(exts, ", ")
}

func generateFileName(original string) string {
	ext := filepath.Ext(original)
	if ext == "" {
		ext = ".bin"
	}
	return uuid.NewString() + ext
}
