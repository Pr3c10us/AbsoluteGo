package middlewares

import (
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/packages/appError"
	"github.com/Pr3c10us/absolutego/packages/response"
	"github.com/Pr3c10us/absolutego/packages/validator"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			var (
				pqErr           *pq.Error
				customError     *appError.CustomError
				validationError *validator.ValidationError
			)
			fmt.Println("error", "Error handler message", zap.Error(err))

			switch {
			case errors.As(err.Err, &pqErr):
				{
					log.Print(pqErr.Code.Name())
					switch pqErr.Code {
					case "23505":
						response.ErrorResponse{
							StatusCode:   http.StatusConflict,
							Message:      "unique key value violated",
							ErrorMessage: pqErr.Detail,
						}.Send(c)
						return
					case "22P02":
						response.ErrorResponse{
							StatusCode:   http.StatusBadRequest,
							Message:      "invalid argument syntax",
							ErrorMessage: pqErr.Message,
						}.Send(c)
						return
					case "23503":
						response.ErrorResponse{
							StatusCode:   http.StatusBadRequest,
							Message:      "invalid foreign key identifier",
							ErrorMessage: pqErr.Detail,
						}.Send(c)
						return
					default:
						response.NewErrorResponse(pqErr).Send(c)
						return
					}
				}
			case errors.As(err.Err, &customError):
				response.ErrorResponse{
					StatusCode:   customError.StatusCode,
					Message:      customError.Message,
					ErrorMessage: customError.ErrorMessage,
				}.Send(c)
				return
			case errors.As(err.Err, &validationError):
				response.ErrorResponse{
					StatusCode:   validationError.StatusCode,
					Message:      validationError.Message,
					ErrorMessage: validationError.ErrorMessage,
				}.Send(c)
			default:
				response.NewErrorResponse(err).Send(c)
				return
			}
		}
	}
}
