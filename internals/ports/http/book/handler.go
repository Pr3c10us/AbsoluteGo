package book

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/services/book"
	"github.com/Pr3c10us/absolutego/internals/services/book/commands"
	"github.com/Pr3c10us/absolutego/packages/appError"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/response"
	"github.com/Pr3c10us/absolutego/packages/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
	"strconv"
	"strings"
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

func (handler *Handler) AddBook(context *gin.Context) {
	var params struct {
		Book    *multipart.FileHeader `form:"book" binding:"required"`
		Chapter string                `form:"chapter" binding:"required,min=1"`
	}

	if err := context.Bind(&params); err != nil {
		err = validator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}

	chapter, err := strconv.Atoi(params.Chapter)
	if err != nil {
		_ = context.Error(err)
		return
	}

	var rootPath = configs.GetRootPath()

	destination, err := generateFileName(params.Book.Filename)
	if err != nil {
		_ = context.Error(err)
		return
	}

	destination = rootPath + "/upload/" + destination

	err = context.SaveUploadedFile(params.Book, destination)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = handler.service.Commands.AddBook.Handle(commands.Parameter{
		File: destination, Chapter: chapter,
	})
	if err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse("Books added", nil, nil).Send(context)
}

func generateFileName(fileName string) (string, error) {
	fileNameParts := strings.Split(fileName, ".")
	if len(fileNameParts) < 2 {
		return "", appError.BadRequest(errors.New("invalid file"))
	}
	return uuid.NewString() + "." + fileNameParts[len(fileNameParts)-1], nil
}
