package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	"net/http"
	"time"
)

func (h *Handler) initWarehouseRoutes(api *gin.RouterGroup) {
	wh := api.Group("/warehouse")
	{
		wh.GET("/:id", h.userIdentity, h.getWarehouse)
		wh.GET("/:id/income-history", h.userIdentity, h.getIncomeHistory)
		wh.GET("/", h.userIdentity, h.getWarehouses)

		// only admin can create, update, delete warehouse
		wh.POST("/", h.adminIdentity, h.createWarehouse)
		wh.PUT("/:id", h.adminIdentity, h.updateWarehouse)
		wh.DELETE("/:id", h.adminIdentity, h.deleteWarehouse)
		wh.GET("/responsible-person", h.adminIdentity, h.getResponsiblePerson)

		wh.GET("/report/:id/xls", h.userIdentity, h.getWarehouseInfoReportXls)
		wh.GET("/report/:id/pdf", h.userIdentity, h.getWarehouseInfoReportPdf)
		wh.GET("report/list/xls", h.userIdentity, h.getWarehouseListReport)
	}
}

// @Summary Get warehouse by id
// @Security ApiKeyAuth
// @Tags warehouse
// @Description Получение склада по id
// @ID get-warehouse
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /warehouse/{id} [GET]
func (h *Handler) getWarehouse(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	wh, err := h.services.Warehouse.GetById(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrNotAllowed) {
			newErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}

		if errors.Is(err, domain.ErrWarehouseNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       wh,
		TotalCount: 1,
	})
}

// @Summary Get warehouses
// @Security ApiKeyAuth
// @Tags warehouse
// @Description Получение списка складов
// @ID get-warehouses
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, name, address, responsible_person, phone, email, max_capacity, current_occupancy, country, region, created_at) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /warehouse [GET]
func (h *Handler) getWarehouses(c *gin.Context) {
	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
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

	whs, count, err := h.services.Warehouse.GetListByCompanyId(c, info.CompanyId, domain.Param{
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
		SortField: field,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       whs,
		TotalCount: count,
	})
}

// @Summary Create warehouse
// @Security ApiKeyAuth
// @Tags warehouse
// @Description Создание склада
// @ID create-warehouse
// @Accept json
// @Produce json
// @Param input body domain.InputWarehouse true "Необходимо указать данные склада."
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /warehouse [POST]
func (h *Handler) createWarehouse(c *gin.Context) {
	var inp domain.InputWarehouse
	if err := c.ShouldBindJSON(&inp); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.services.Warehouse.Create(c, domain.Warehouse{
		Name:              inp.Name,
		Address:           inp.Address,
		ResponsiblePerson: inp.ResponsiblePerson,
		Phone:             inp.Phone,
		Email:             inp.Email,
		MaxCapacity:       inp.MaxCapacity,
		CurrentOccupancy:  inp.CurrentOccupancy,
		OtherFields:       inp.OtherFields,
		Country:           inp.Country,
		Region:            inp.Region,
		Locality:          inp.Locality,
		Comments:          inp.Comments,
		CompanyId:         info.CompanyId,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newCreateSuccessIdResponse(c, id)
}

// @Summary Update warehouse
// @Security ApiKeyAuth
// @Tags warehouse
// @Description Обновление склада своей компании
// @ID update-warehouse
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Param input body domain.InputWarehouse true "Необходимо указать данные склада."
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /warehouse/{id} [PUT]
func (h *Handler) updateWarehouse(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var inp domain.WarehouseUpdate
	if err := c.ShouldBindJSON(&inp); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	inp.ID = id

	if err := h.services.Warehouse.Update(c, inp, info); err != nil {
		if errors.Is(err, domain.ErrWarehouseNotFound) {
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

// @Summary Delete warehouse
// @Security ApiKeyAuth
// @Tags warehouse
// @Description Удаление склада своей компании
// @ID delete-warehouse
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /warehouse/{id} [DELETE]
func (h *Handler) deleteWarehouse(c *gin.Context) {
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

	if err = h.services.Warehouse.Delete(c, id, info); err != nil {
		if errors.Is(err, domain.ErrWarehouseNotFound) {
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

// @Summary Get responsible users
// @Security ApiKeyAuth
// @Tags warehouse
// @Description Получение списка доступных ответственных лиц для склада
// @ID get-responsible-person
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, username, name, email, phone, created_at, updated_at, country, position) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /warehouse/responsible-person [GET]
func (h *Handler) getResponsiblePerson(c *gin.Context) {
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

	users, count, err := h.services.Warehouse.GetResponsibleUsers(c, info.CompanyId, domain.Param{
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
		SortField: field,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       users,
		TotalCount: count,
	})
}

// @Summary Get warehouse income history
// @Security ApiKeyAuth
// @Tags warehouse
// @Description Получение истории поступлений товаров
// @ID get-income-history
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(received_date) default(received_date)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /warehouse/{id}/income-history [GET]
func (h *Handler) getIncomeHistory(c *gin.Context) {
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

	history, count, err := h.services.Warehouse.GetIncomeHistoryByWarehouseId(c, id, domain.Param{
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
		Data:       history,
		TotalCount: count,
	})
}

// @Summary      Generate warehouse info xls
// @Security 	 ApiKeyAuth
// @Tags         warehouse
// @Description  Возвращает XLS файл с данными склада
// @ID           warehouse-info-report-xls
// @Accept       json
// @Produce      application/octet-stream
// @Param 	  	 id path int true "Warehouse ID"
// @Success      200 {file} file "XLS файл отчета"
// @Failure 	 400,404 {object} domain.ErrorResponse
// @Failure 	 500 {object} domain.ErrorResponse
// @Failure 	 default {object} domain.ErrorResponse
// @Router       /warehouse/report/{id}/xls [GET]
func (h *Handler) getWarehouseInfoReportXls(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	report, err := h.services.Warehouse.GenerateWarehouseInfoReportXls(c, id, info)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	currentTime := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("warehouse_%d_report_%s.xlsx", id, currentTime)

	// Устанавливаем заголовки и отправляем файл
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))

	if err = report.Write(c.Writer); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
}

// @Summary      Generate warehouse info pdf
// @Security 	 ApiKeyAuth
// @Tags         warehouse
// @Description  Возвращает PDF файл с данными склада
// @ID           warehouse-info-report-pdf
// @Accept       json
// @Produce      application/pdf
// @Param 		 id path int true "Warehouse ID"
// @Success      200 {file} file "PDF файл отчета"
// @Failure 	 400,404 {object} domain.ErrorResponse
// @Failure 	 500 {object} domain.ErrorResponse
// @Failure 	 default {object} domain.ErrorResponse
// @Router       /warehouse/report/{id}/pdf [GET]
func (h *Handler) getWarehouseInfoReportPdf(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	report, err := h.services.Warehouse.GenerateWarehouseInfoReportPdf(c, id, info)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Формируем название файла с текущей датой и временем
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("warehouse_%d_report_%s.pdf", id, currentTime)

	// Устанавливаем заголовки и отправляем файл
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	if err = report.Output(c.Writer); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) getWarehouseListReport(c *gin.Context) {

}
