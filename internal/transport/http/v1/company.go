package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	"net/http"
	"time"
)

func (h *Handler) initCompanyRoutes(api *gin.RouterGroup) {
	company := api.Group("/company")
	{
		company.GET("/info", h.userIdentity, h.getCompanyInfo)

		// only admin can update company
		company.PUT("/:id", h.adminIdentity, h.updateCompany)

		// only super admin can create & delete company
		company.POST("/", h.superAdminIdentity, h.createCompany)
		company.GET("/:id", h.superAdminIdentity, h.getCompany)
		company.DELETE("/:id", h.superAdminIdentity, h.deleteCompany)
		company.GET("/", h.superAdminIdentity, h.getCompanies)
	}
}

// @Summary Get company info
// @Security ApiKeyAuth
// @Tags company
// @Description Получение информации о компании для пользователя
// @ID get-company-info
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /company/info [GET]
func (h *Handler) getCompanyInfo(c *gin.Context) {
	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	companyInfo, err := h.services.Company.GetById(c.Request.Context(), info.CompanyId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       companyInfo,
		TotalCount: 1,
	})
}

// @Summary Update company
// @Security ApiKeyAuth
// @Tags company
// @Description Обновление компании.
// @Description Только super admin может обновлять active & approve компании
// @Description Для обновления указывать только необходимые поля.
// @ID update-company
// @Accept  json
// @Produce  json
// @Param id path int true "Company ID" example(1)
// @Param company body domain.CompanyUpdate true "Company info"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /company/{id} [PUT]
func (h *Handler) updateCompany(c *gin.Context) {
	var req domain.CompanyUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	req.ID, err = parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = h.services.Company.Update(c.Request.Context(), req, info); err != nil {
		if errors.Is(err, domain.ErrInvalidTimezone) {
			newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if errors.Is(err, domain.ErrCompanyNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, domain.ErrNotAllowed) {
			newErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Create company
// @Security ApiKeyAuth
// @Tags company
// @Description Создание компании.
// @Description Только super admin может создавать компании
// @ID create-company
// @Accept  json
// @Produce  json
// @Param company body domain.CreateCompany true "Company info"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /company/ [POST]
func (h *Handler) createCompany(c *gin.Context) {
	var req domain.CreateCompany
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Company.Create(c.Request.Context(), domain.Company{
		NameRu:     req.NameRu,
		NameEn:     req.NameEn,
		Country:    req.Country,
		Address:    req.Address,
		Phone:      req.Phone,
		Email:      req.Email,
		Website:    req.Website,
		IsActive:   req.IsActive,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		IsApproved: req.IsApproved,
		Timezone:   req.Timezone,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newCreateSuccessIdResponse(c, id)
}

// @Summary Delete company
// @Security ApiKeyAuth
// @Tags company
// @Description Удаление компании.
// @Description Только super admin может удалять компании.
// @ID delete-company
// @Accept  json
// @Produce  json
// @Param id path int true "Company ID"
// @Success 200 {object} domain.MessageResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /company/{id} [DELETE]
func (h *Handler) deleteCompany(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := h.services.Company.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Get company
// @Security ApiKeyAuth
// @Tags company
// @Description Получение информации о компании по ID.
// @Description Только super admin может получать информацию по id компании.
// @ID get-company
// @Accept  json
// @Produce  json
// @Param id path int true "Company ID" example(1)
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /company/{id} [GET]
func (h *Handler) getCompany(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	companyInfo, err := h.services.Company.GetById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       companyInfo,
		TotalCount: 1,
	})
}

// @Summary Get companies
// @Security ApiKeyAuth
// @Tags company
// @Description Получение информации по всем компаниям.
// @Description Только super admin может получать информацию.
// @ID get-companies
// @Accept  json
// @Produce  json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, name_ru, name_en, country, address, phone, email, website, is_active, created_at, updated_at, is_approved) default(name_ru)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /company/ [GET]
func (h *Handler) getCompanies(c *gin.Context) {
	limit, err := parseLimitQueryParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	offset, err := parseOffsetQueryParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	list, count, err := h.services.Company.List(c.Request.Context(), domain.Param{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       list,
		TotalCount: count,
	})
}
