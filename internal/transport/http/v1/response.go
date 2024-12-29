package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/rusystem/crm-api/pkg/logger"
	"net/http"
)

func newSuccessResponse(c *gin.Context, code int, resp domain.SuccessResponse) {
	c.JSON(code, resp)
}

func newErrorResponse(c *gin.Context, code int, message string) {
	logger.Error(message)

	c.AbortWithStatusJSON(code, domain.ErrorResponse{
		Code:    code,
		IsError: true,
		Message: message,
	})
}

func newSuccessOkResponse(c *gin.Context) {
	c.JSON(http.StatusOK, domain.SuccessResponse{
		Data: domain.MessageResponse{Message: "success"},
	})
}

func newCreateSuccessIdResponse(c *gin.Context, id int64) {
	newSuccessResponse(c, http.StatusCreated, domain.SuccessResponse{
		Data:       domain.IdResponse{ID: id},
		TotalCount: 1,
	})
}
