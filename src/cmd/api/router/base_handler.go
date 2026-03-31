package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Result bool        `json:"result" example:"true"`
	Data   interface{} `json:"data"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Result bool        `json:"result" example:"false"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error" example:"error message"`
}

// OK returns success response {result: true, data: ...} with 200 status code
func (handler *BaseHandler) OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Result: true,
		Data:   data,
	})
}

// FailedInternal returns 500 error {result: false, data: ..., error: ...}
func (handler *BaseHandler) FailedInternal(c *gin.Context, data interface{}, err error) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Result: false,
		Data:   data,
		Error:  err.Error(),
	})
}

// FailedBadRequest returns 400 error {result: false, data: ..., error: ...}
func (handler *BaseHandler) FailedBadRequest(c *gin.Context, data interface{}, err error) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Result: false,
		Data:   data,
		Error:  err.Error(),
	})
}
