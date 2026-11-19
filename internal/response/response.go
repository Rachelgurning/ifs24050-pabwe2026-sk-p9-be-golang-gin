package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string { return e.Message }

func NewError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Success(c *gin.Context, code int, message string, data any) {
	body := gin.H{"status": "success", "message": message}
	if data != nil {
		body["data"] = data
	}
	c.JSON(code, body)
}

func Fail(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"status": "fail", "message": message})
}

func Error(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"status": "error", "message": message})
}

func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.Code >= 500 {
			Error(c, appErr.Code, appErr.Message)
		} else {
			Fail(c, appErr.Code, appErr.Message)
		}
		return
	}
	Error(c, http.StatusInternalServerError, "Terjadi kesalahan pada server")
}
