package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	"net/http"
)

func (h *Handler) initUnitOfMeasureRoutes(api *gin.RouterGroup) {
	measure := api.Group("/measure", h.userIdentity)
	{
		measure.POST("/", h.createMeasure)
		measure.GET("/:id", h.getMeasureById)
		measure.PUT("/:id", h.updateMeasure)
		measure.DELETE("/:id", h.deleteMeasure)
		measure.GET("/", h.getMeasureList)
	}
}

// @Summary Create unit of measure
// @Security ApiKeyAuth
// @Tags unit of measure
// @Description Создание единицы измерения
// @ID create-unit-of-measure
// @Accept json
// @Produce json
// @Param input body domain.CreateUnitOfMeasure true "Необходимо указать данные единицы измерения"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /measure [POST]
func (h *Handler) createMeasure(c *gin.Context) {
	var inp domain.CreateUnitOfMeasure
	if err := c.ShouldBindJSON(&inp); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.services.UnitOfMeasure.Create(c.Request.Context(), domain.UnitOfMeasure{
		Name:         inp.Name,
		NameEn:       inp.NameEn,
		Abbreviation: inp.Abbreviation,
		Description:  inp.Description,
		CompanyID:    info.CompanyId,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newCreateSuccessIdResponse(c, id)
}

// @Summary Get unit of measure
// @Security ApiKeyAuth
// @Tags unit of measure
// @Description Получение единицы измерения
// @ID get-unit-of-measure
// @Accept json
// @Produce json
// @Param id path int true "ID единицы измерения"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /measure/{id} [GET]
func (h *Handler) getMeasureById(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	measure, err := h.services.UnitOfMeasure.GetById(c.Request.Context(), id, info.CompanyId)
	if err != nil {
		if errors.Is(err, domain.ErrUnitOfMeasureNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       measure,
		TotalCount: 1,
	})
}

// @Summary Update unit of measure
// @Security ApiKeyAuth
// @Tags unit of measure
// @Description Обновление единицы измерения
// @ID update-unit-of-measure
// @Accept json
// @Produce json
// @Param id path int true "ID единицы измерения"
// @Param input body domain.UpdateUnitOfMeasure true "Необходимо указать данные единицы измерения"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /measure/{id} [PUT]
func (h *Handler) updateMeasure(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var inp domain.UpdateUnitOfMeasure
	if err = c.ShouldBindJSON(&inp); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	inp.ID = id
	inp.CompanyID = info.CompanyId

	if err = h.services.UnitOfMeasure.Update(c.Request.Context(), inp); err != nil {
		if errors.Is(err, domain.ErrUnitOfMeasureNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Delete unit of measure
// @Security ApiKeyAuth
// @Tags unit of measure
// @Description Удаление единицы измерения
// @ID delete-unit-of-measure
// @Accept json
// @Produce json
// @Param id path int true "ID единицы измерения"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /measure/{id} [DELETE]
func (h *Handler) deleteMeasure(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err = h.services.UnitOfMeasure.Delete(c, id, info.CompanyId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Get unit of measure list
// @Security ApiKeyAuth
// @Tags unit of measure
// @Description Список единиц измерений
// @ID get-unit-of-measure-list
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, name, name_en, abbreviation, description) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /measure [GET]
func (h *Handler) getMeasureList(c *gin.Context) {
	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	sort, field, err := parseSortParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

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

	measures, count, err := h.services.UnitOfMeasure.List(c.Request.Context(), domain.Param{
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
		SortField: field,
		CompanyId: info.CompanyId,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       measures,
		TotalCount: count,
	})
}
