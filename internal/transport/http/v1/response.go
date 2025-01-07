package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/rusystem/crm-api/pkg/logger"
	"net/http"
)

var validationMessages = map[string]map[string]string{
	"Username": {
		"required": "Username is required",
		"min":      "Username must be at least 5 characters long",
		"max":      "Username must not exceed 140 characters",
	},
	"Name": {
		"required": "Name is required",
		"max":      "Name must not exceed 140 characters",
	},
	"Email": {
		"required": "Email is required",
		"email":    "Email must be a valid email address",
		"min":      "Email must be at least 5 characters long",
		"max":      "Email must not exceed 140 characters",
	},
	"Phone": {
		"required": "Phone number is required",
		"min":      "Phone number must be at least 7 characters long",
		"max":      "Phone number must not exceed 140 characters",
	},
	"Password": {
		"required": "Password is required",
		"min":      "Password must be at least 8 characters long",
		"max":      "Password must not exceed 255 characters",
	},
	"NameRu": {
		"required": "NameRu (company name in Russian) is required",
	},
	"Address": {
		"required": "Address is required",
		"min":      "Password must be at least 5 characters long",
		"max":      "Password must not exceed 140 characters",
	},
	"LegalAddress": {
		"required": "Legal address is required",
	},
	"ActualAddress": {
		"required": "Actual address is required",
	},
	"WarehouseAddress": {
		"required": "Warehouse address is required",
	},
	"ContactPerson": {
		"required": "Contact person is required",
		"max":      "Contact person must not exceed 140 characters",
	},
	"ContractNumber": {
		"required": "Contract number is required",
	},
	"ContractDate": {
		"required": "Contract date is required",
	},
	"Balance": {
		"required": "Balance is required",
	},
	"BankDetails": {
		"required": "Bank details are required",
	},
	"ID": {
		"required": "User ID is required",
	},
	"ResponsiblePerson": {
		"required": "Responsible person is required",
		"min":      "Responsible person must be at least 5 characters long",
		"max":      "Responsible person must not exceed 140 characters",
	},
}

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

func newBindingErrorResponse(c *gin.Context, err error) {
	logger.Error(err.Error())

	code, message := validationErrorHandler(err)

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

// ValidationErrorHandler обёртка для обработки ошибок валидации
func validationErrorHandler(err error) (int, string) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		// Извлекаем первую ошибку
		firstError := validationErrors[0]
		field := firstError.Field()
		tag := firstError.Tag()

		// Формируем пользовательское сообщение
		message := generateErrorMessage(field, tag)

		return http.StatusUnprocessableEntity, message
	}

	// Обработка других ошибок
	return http.StatusBadRequest, "invalid input data"
}

// generateErrorMessage создает сообщение об ошибке на основе карты
func generateErrorMessage(field, tag string) string {
	if fieldMessages, ok := validationMessages[field]; ok {
		if message, ok := fieldMessages[tag]; ok {
			return message
		}
	}
	// Возврат общего сообщения, если специфическое не найдено
	switch tag {
	case "required":
		return field + " is required"
	case "min":
		return field + " is too short"
	case "max":
		return field + " is too long"
	default:
		return "Invalid value for " + field
	}
}
