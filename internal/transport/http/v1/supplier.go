package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	"net/http"
)

func (h *Handler) initSupplierRoutes(api *gin.RouterGroup) {
	spl := api.Group("/supplier")
	{
		spl.GET("/:id", h.userIdentity, h.getSupplier)
		spl.GET("/", h.userIdentity, h.getSuppliers)

		// only admin can create, update, delete supplier
		spl.POST("/", h.adminIdentity, h.createSupplier)
		spl.PUT("/:id", h.adminIdentity, h.updateSupplier)
		spl.DELETE("/:id", h.adminIdentity, h.deleteSupplier)
	}
}

// @Summary Get supplier by id
// @Security ApiKeyAuth
// @Tags supplier
// @Description Получение поставщика по id
// @ID get-supplier
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /supplier/{id} [GET]
func (h *Handler) getSupplier(c *gin.Context) {
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

	spl, err := h.services.Supplier.GetById(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrSupplierNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       spl,
		TotalCount: 1,
	})
}

// @Summary Get suppliers
// @Security ApiKeyAuth
// @Tags supplier
// @Description Получение всех поставщиков компании
// @ID get-suppliers
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, name, legal_address, actual_address, warehouse_address, contact_person, phone, email, website, contract_number, product_categories, purchase_amount, balance, product_types, country, region, tax_id, registration_date, is_active) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /supplier [GET]
func (h *Handler) getSuppliers(c *gin.Context) {
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
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	offset, err := parseOffsetQueryParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	spls, count, err := h.services.Supplier.GetListByCompanyId(c, info.CompanyId, domain.Param{
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
		Data:       spls,
		TotalCount: count,
	})
}

// @Summary Create supplier
// @Security ApiKeyAuth
// @Tags supplier
// @Description Создание поставщика
// @ID create-supplier
// @Accept json
// @Produce json
// @Param input body domain.InputSupplier true "Необходимо указать данные поставщика."
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /supplier [POST]
func (h *Handler) createSupplier(c *gin.Context) {
	var inp domain.InputSupplier
	if err := c.ShouldBindJSON(&inp); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.services.Supplier.Create(c, domain.Supplier{
		Name:              inp.Name,
		LegalAddress:      inp.LegalAddress,
		ActualAddress:     inp.ActualAddress,
		WarehouseAddress:  inp.WarehouseAddress,
		ContactPerson:     inp.ContactPerson,
		Phone:             inp.Phone,
		Email:             inp.Email,
		Website:           inp.Website,
		ContractNumber:    inp.ContractNumber,
		ProductCategories: inp.ProductCategories,
		PurchaseAmount:    inp.PurchaseAmount,
		Balance:           inp.Balance,
		ProductTypes:      inp.ProductTypes,
		Comments:          inp.Comments,
		Files:             inp.Files,
		Country:           inp.Country,
		Region:            inp.Region,
		TaxID:             inp.TaxID,
		BankDetails:       inp.BankDetails,
		RegistrationDate:  inp.RegistrationDate,
		PaymentTerms:      inp.PaymentTerms,
		IsActive:          inp.IsActive,
		OtherFields:       inp.OtherFields,
		CompanyId:         info.CompanyId,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newCreateSuccessIdResponse(c, id)
}

// @Summary Delete supplier
// @Security ApiKeyAuth
// @Tags supplier
// @Description Удаление поставщика своей компании
// @ID delete-supplier
// @Accept json
// @Produce json
// @Param id path int true "ID поставщика"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /supplier/{id} [DELETE]
func (h *Handler) deleteSupplier(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err = h.services.Supplier.Delete(c, id, info); err != nil {
		if errors.Is(err, domain.ErrSupplierNotFound) {
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

// @Summary Update supplier
// @Security ApiKeyAuth
// @Tags supplier
// @Description Обновление поставщика
// @ID update-supplier
// @Accept json
// @Produce json
// @Param id path int true "ID поставщика"
// @Param input body domain.UpdateSupplier true "Необходимо указать данные поставщика."
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /supplier/{id} [PUT]
func (h *Handler) updateSupplier(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var inp domain.UpdateSupplier
	if err := c.ShouldBindJSON(&inp); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	inp.Id = id

	if err = h.services.Supplier.Update(c, inp, info); err != nil {
		if errors.Is(err, domain.ErrSupplierNotFound) {
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
