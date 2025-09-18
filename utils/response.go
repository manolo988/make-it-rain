package utils

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Error: message,
	})
}

func RespondWithErrorDetails(c *gin.Context, code int, message, errorCode, details string) {
	c.JSON(code, ErrorResponse{
		Error:   message,
		Code:    errorCode,
		Details: details,
	})
}

func RespondWithSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, SuccessResponse{
		Message: message,
		Data:    data,
	})
}