package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/rusystem/crm-api/tools"
	"net/http"
	"time"
)

func (h *Handler) initMaterialsRoutes(api *gin.RouterGroup) {
	materials := api.Group("/materials", h.userIdentity) //todo добавить секции в мидлвейер
	{
		planning := materials.Group("/planning")
		{
			planning.POST("/", h.createPlanning)
			planning.GET("/:id", h.getPlanningById)
			planning.PUT("/:id", h.updatePlanningById)
			planning.DELETE("/:id", h.deletePlanningById)
			planning.GET("/", h.getPlanningList)
			planning.PUT("/move-to-purchased/:id", h.movePlanningToPurchased)
		}

		purchased := materials.Group("/purchased")
		{
			purchased.POST("/", h.createPurchased)
			purchased.GET("/:id", h.getPurchasedById)
			purchased.PUT("/:id", h.updatePurchasedById)
			purchased.DELETE("/:id", h.deletePurchasedById)
			purchased.GET("/", h.getPurchasedList)
			purchased.GET("/:id/qr-code", h.getPurchasedQrCode)
			purchased.GET("/:id/barcode", h.getPurchasedBarcode)
			purchased.PUT("/move-to-archive/:id", h.movePurchasedToArchive)
		}

		archive := materials.Group("/archive")
		{
			planning := archive.Group("/planning")
			{
				planning.GET("/:id", h.getPlanningArchiveById)
				planning.GET("/", h.getPlanningArchiveList)
				planning.DELETE("/:id", h.deletePlanningArchiveById)
			}

			purchased := archive.Group("/purchased")
			{
				purchased.GET("/:id", h.getPurchasedArchiveById)
				purchased.GET("/", h.getPurchasedArchiveList)
				purchased.DELETE("/:id", h.deletePurchasedArchiveById)
			}
		}

		search := materials.Group("/search")
		{
			search.GET("/", h.searchMaterial)
		}

		category := materials.Group("/category")
		{
			category.POST("/", h.createCategory)
			category.GET("/:id", h.getCategoryById)
			category.PUT("/:id", h.updateCategory)
			category.DELETE("/:id", h.deleteCategory)
			category.GET("/", h.getCategoryList)
			category.GET("/search", h.searchCategory)
		}
	}
}

// @Summary Create planning material
// @Security ApiKeyAuth
// @Tags materials planning
// @Description Создание планируемого материала
// @ID create-planning-material
// @Accept json
// @Produce json
// @Param input body domain.CreatePlanningMaterial true "Необходимо указать данные планируемого материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/planning [POST]
func (h *Handler) createPlanning(c *gin.Context) {
	var inp domain.CreatePlanningMaterial
	if err := c.ShouldBindJSON(&inp); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.services.Materials.CreatePlanning(c, info, domain.Material{
		WarehouseID:            inp.WarehouseID,
		Name:                   inp.Name,
		ByInvoice:              inp.ByInvoice,
		Article:                inp.Article,
		ProductCategory:        inp.ProductCategory,
		Unit:                   inp.Unit,
		TotalQuantity:          inp.TotalQuantity,
		Volume:                 inp.Volume,
		PriceWithoutVAT:        inp.PriceWithoutVAT,
		TotalWithoutVAT:        inp.TotalWithoutVAT,
		SupplierID:             inp.SupplierID,
		ContractDate:           inp.ContractDate,
		File:                   inp.File,
		Status:                 inp.Status,
		Comments:               inp.Comments,
		Reserve:                inp.Reserve,
		ReceivedDate:           inp.ReceivedDate,
		LastUpdated:            time.Now().UTC(),
		MinStockLevel:          inp.MinStockLevel,
		ExpirationDate:         inp.ExpirationDate,
		ResponsiblePerson:      inp.ResponsiblePerson,
		StorageCost:            inp.StorageCost,
		WarehouseSection:       inp.WarehouseSection,
		IncomingDeliveryNumber: inp.IncomingDeliveryNumber,
		OtherFields:            inp.OtherFields,
		CompanyID:              info.CompanyId,
		InternalName:           inp.InternalName,
		UnitsPerPackage:        inp.UnitsPerPackage,
		ContractNumber:         inp.ContractNumber,
	})
	if err != nil {
		if errors.Is(err, domain.ErrWarehouseNotFound) || errors.Is(err, domain.ErrSupplierNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newCreateSuccessIdResponse(c, id)
}

// @Summary Get planning material
// @Security ApiKeyAuth
// @Tags materials planning
// @Description Получение планируемого материала
// @ID get-planning-material
// @Accept json
// @Produce json
// @Param id path int true "ID планируемого материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/planning/{id} [GET]
func (h *Handler) getPlanningById(c *gin.Context) {
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

	material, err := h.services.Materials.GetPlanningById(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       material,
		TotalCount: 1,
	})
}

// @Summary Update planning material
// @Security ApiKeyAuth
// @Tags materials planning
// @Description Обновление планируемого материала
// @ID update-planning-material
// @Accept json
// @Produce json
// @Param id path int true "ID планируемого материала"
// @Param input body domain.UpdatePlanningMaterial true "Необходимо указать данные планируемого материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/planning/{id} [PUT]
func (h *Handler) updatePlanningById(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var inp domain.UpdatePlanningMaterial
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

	if err = h.services.Materials.UpdatePlanningById(c, inp, info); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, domain.ErrNotAllowed) {
			newErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}

		if errors.Is(err, domain.ErrWarehouseNotFound) || errors.Is(err, domain.ErrSupplierNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Delete planning material
// @Security ApiKeyAuth
// @Tags materials planning
// @Description Удаление планируемого материала
// @ID delete-planning-material
// @Accept json
// @Produce json
// @Param id path int true "ID планируемого материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/planning/{id} [DELETE]
func (h *Handler) deletePlanningById(c *gin.Context) {
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

	if err = h.services.Materials.DeletePlanningById(c, id, info); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

// @Summary Get planning list
// @Security ApiKeyAuth
// @Tags materials planning
// @Description Список планируемых материалов
// @ID get-planning-list
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, warehouse_id, item_id, name, article, product_category, total_quantity, volume, price_without_vat, total_without_vat, supplier_id, location, status, received_date, last_updated, min_stock_level, expiration_date, storage_cost, warehouse_section, incoming_delivery_number) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/planning [GET]
func (h *Handler) getPlanningList(c *gin.Context) {
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

	mtrls, count, err := h.services.Materials.GetPlanningList(c.Request.Context(), domain.MaterialParams{
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
		Data:       mtrls,
		TotalCount: count,
	})
}

// @Summary Move planning material to purchased
// @Security ApiKeyAuth
// @Tags materials planning
// @Description Перемещение планируемого материала в закупленные
// @ID move-planning-to-purchased
// @Accept json
// @Produce json
// @Param id path int true "ID планируемого материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/planning/move-to-purchased/{id} [PUT]
func (h *Handler) movePlanningToPurchased(c *gin.Context) {
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

	newId, itemId, err := h.services.Materials.MovePlanningToPurchased(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data: domain.PurchasedIdResponse{ID: newId, ItemId: itemId},
	})
}

// @Summary Create purchased material
// @Security ApiKeyAuth
// @Tags materials purchased
// @Description Создание закупленного материала
// @ID create-purchased-material
// @Accept json
// @Produce json
// @Param input body domain.CreatePurchasedMaterial true "Необходимо указать данные закупленного материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/purchased [POST]
func (h *Handler) createPurchased(c *gin.Context) {
	var inp domain.CreatePurchasedMaterial
	if err := c.ShouldBindJSON(&inp); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, itemId, err := h.services.Materials.CreatePurchased(c, info, domain.Material{
		WarehouseID:            inp.WarehouseID,
		Name:                   inp.Name,
		ByInvoice:              inp.ByInvoice,
		Article:                inp.Article,
		ProductCategory:        inp.ProductCategory,
		Unit:                   inp.Unit,
		TotalQuantity:          inp.TotalQuantity,
		Volume:                 inp.Volume,
		PriceWithoutVAT:        inp.PriceWithoutVAT,
		TotalWithoutVAT:        inp.TotalWithoutVAT,
		SupplierID:             inp.SupplierID,
		Location:               inp.Location,
		ContractDate:           inp.ContractDate,
		File:                   inp.File,
		Status:                 inp.Status,
		Comments:               inp.Comments,
		Reserve:                inp.Reserve,
		ReceivedDate:           inp.ReceivedDate,
		LastUpdated:            time.Now().UTC(),
		MinStockLevel:          inp.MinStockLevel,
		ExpirationDate:         inp.ExpirationDate,
		ResponsiblePerson:      inp.ResponsiblePerson,
		StorageCost:            inp.StorageCost,
		WarehouseSection:       inp.WarehouseSection,
		IncomingDeliveryNumber: inp.IncomingDeliveryNumber,
		OtherFields:            inp.OtherFields,
		CompanyID:              info.CompanyId,
		InternalName:           inp.InternalName,
		UnitsPerPackage:        inp.UnitsPerPackage,
		ContractNumber:         inp.ContractNumber,
	})
	if err != nil {
		if errors.Is(err, domain.ErrWarehouseNotFound) || errors.Is(err, domain.ErrSupplierNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusCreated, domain.SuccessResponse{
		Data: domain.PurchasedIdResponse{ID: id, ItemId: itemId},
	})
}

// @Summary Get purchased material
// @Security ApiKeyAuth
// @Tags materials purchased
// @Description Получение закупленного материала
// @ID get-purchased-material
// @Accept json
// @Produce json
// @Param id path int true "ID закупленного материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/purchased/{id} [GET]
func (h *Handler) getPurchasedById(c *gin.Context) {
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

	material, err := h.services.Materials.GetPurchasedById(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       material,
		TotalCount: 1,
	})
}

// @Summary Update purchased material
// @Security ApiKeyAuth
// @Tags materials purchased
// @Description Обновление закупленного материала
// @ID update-purchased-material
// @Accept json
// @Produce json
// @Param id path int true "ID закупленного материала"
// @Param input body domain.UpdatePurchasedMaterial true "Необходимо указать данные закупленного материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/purchased/{id} [PUT]
func (h *Handler) updatePurchasedById(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var inp domain.UpdatePurchasedMaterial
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

	if err = h.services.Materials.UpdatePurchasedById(c, inp, info); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, domain.ErrNotAllowed) {
			newErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}

		if errors.Is(err, domain.ErrWarehouseNotFound) || errors.Is(err, domain.ErrSupplierNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Delete purchased material
// @Security ApiKeyAuth
// @Tags materials purchased
// @Description Удаление закупленного материала
// @ID delete-purchased-material
// @Accept json
// @Produce json
// @Param id path int true "ID закупленного материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/purchased/{id} [DELETE]
func (h *Handler) deletePurchasedById(c *gin.Context) {
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

	if err = h.services.Materials.DeletePurchasedById(c, id, info); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

// @Summary Get purchased list
// @Security ApiKeyAuth
// @Tags materials purchased
// @Description Получение списка закупленных материалов
// @ID get-purchased-list
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, warehouse_id, item_id, name, article, product_category, total_quantity, volume, price_without_vat, total_without_vat, supplier_id, location, status, received_date, last_updated, min_stock_level, expiration_date, storage_cost, warehouse_section, incoming_delivery_number) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/purchased [GET]
func (h *Handler) getPurchasedList(c *gin.Context) {
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

	mtrls, count, err := h.services.Materials.GetPurchasedList(c.Request.Context(), domain.MaterialParams{
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
		Data:       mtrls,
		TotalCount: count,
	})
}

// @Summary Get purchased QR code
// @Security ApiKeyAuth
// @Tags materials purchased
// @Description Получение QR кода закупленного материала
// @ID get-purchased-qr-code
// @Accept json
// @Produce  image/png
// @Param id path int true "ID закупленного материала"
// @Success 200 {file} png "QR-код в формате PNG"
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/purchased/{id}/qr-code [GET]
func (h *Handler) getPurchasedQrCode(c *gin.Context) {
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

	material, err := h.services.Materials.GetPurchasedById(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

	qrCode, err := tools.GenerateQRCodePNG(domain.CodeInfo{
		Id:     material.ID,
		ItemId: material.ItemID,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "image/png")
	c.Data(http.StatusOK, "image/png", qrCode)
}

// @Summary Get purchased barcode
// @Security ApiKeyAuth
// @Tags materials purchased
// @Description Получение штрихкода закупленного материала
// @ID get-purchased-barcode
// @Accept json
// @Produce  image/png
// @Param id path int true "ID закупленного материала"
// @Success 200 {file} png "Штрихкод"
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/purchased/{id}/barcode [GET]
func (h *Handler) getPurchasedBarcode(c *gin.Context) {
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

	material, err := h.services.Materials.GetPurchasedById(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

	width := 300
	height := 100 //todo возможно вынести в параметры

	barCode, err := tools.GenerateBarcode(domain.CodeInfo{
		Id:     material.ID,
		ItemId: material.ItemID,
	}, width, height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "image/png")
	c.Data(http.StatusOK, "image/png", barCode)
}

// @Summary Move purchased to archive
// @Security ApiKeyAuth
// @Tags materials purchased
// @Description Перемещение закупленного материала в архив
// @ID move-purchased-to-archive
// @Accept json
// @Produce json
// @Param id path int true "ID закупленного материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/purchased/move-to-archive/{id} [PUT]
func (h *Handler) movePurchasedToArchive(c *gin.Context) {
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

	if err = h.services.Materials.MovePurchasedToArchive(c, id, info); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

// @Summary Get planning archive by id
// @Security ApiKeyAuth
// @Tags materials archive
// @Description Получение запланированного материала из архива по ID
// @ID get-planning-archive-by-id
// @Accept json
// @Produce json
// @Param id path int true "ID архиввного планируемого материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/archive/planning/{id} [GET]
func (h *Handler) getPlanningArchiveById(c *gin.Context) {
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

	material, err := h.services.Materials.GetPlanningArchiveById(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       material,
		TotalCount: 1,
	})
}

// @Summary Get planning archive list
// @Security ApiKeyAuth
// @Tags materials archive
// @Description Получение списка запланированных материалов из архива
// @ID get-planning-archive-list
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, warehouse_id, item_id, name, article, product_category, total_quantity, volume, price_without_vat, total_without_vat, supplier_id, location, status, received_date, last_updated, min_stock_level, expiration_date, storage_cost, warehouse_section, incoming_delivery_number) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/archive/planning [GET]
func (h *Handler) getPlanningArchiveList(c *gin.Context) {
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

	mtrls, count, err := h.services.Materials.GetPlanningArchiveList(c.Request.Context(), domain.MaterialParams{
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
		Data:       mtrls,
		TotalCount: count,
	})
}

// @Summary Delete planning archive by id
// @Security ApiKeyAuth
// @Tags materials archive
// @Description Удаление запланированного материала из архива по ID
// @ID delete-planning-archive-by-id
// @Accept json
// @Produce json
// @Param id path int true "ID архиввного планируемого материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/archive/planning/{id} [DELETE]
func (h *Handler) deletePlanningArchiveById(c *gin.Context) {
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

	if err = h.services.Materials.DeletePlanningArchiveById(c, id, info); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

// @Summary Get purchased archive by id
// @Security ApiKeyAuth
// @Tags materials archive
// @Description Получение закупленного материала из архива по ID
// @ID get-purchased-archive-by-id
// @Accept json
// @Produce json
// @Param id path int true "ID архиввного закупленного материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/archive/purchased/{id} [GET]
func (h *Handler) getPurchasedArchiveById(c *gin.Context) {
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

	material, err := h.services.Materials.GetPurchasedArchiveById(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       material,
		TotalCount: 1,
	})
}

// @Summary Get purchased archive list
// @Security ApiKeyAuth
// @Tags materials archive
// @Description Получение списка закупленных материалов из архива
// @ID get-purchased-archive-list
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, warehouse_id, item_id, name, article, product_category, total_quantity, volume, price_without_vat, total_without_vat, supplier_id, location, status, received_date, last_updated, min_stock_level, expiration_date, storage_cost, warehouse_section, incoming_delivery_number) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/archive/purchased [GET]
func (h *Handler) getPurchasedArchiveList(c *gin.Context) {
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

	mtrls, count, err := h.services.Materials.GetPurchasedArchiveList(c.Request.Context(), domain.MaterialParams{
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
		Data:       mtrls,
		TotalCount: count,
	})
}

// @Summary Delete purchased archive by id
// @Security ApiKeyAuth
// @Tags materials archive
// @Description Удаление закупленного материала из архива по ID
// @ID delete-purchased-archive-by-id
// @Accept json
// @Produce json
// @Param id path int true "ID архиввного закупленного материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/archive/purchased/{id} [DELETE]
func (h *Handler) deletePurchasedArchiveById(c *gin.Context) {
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

	if err = h.services.Materials.DeletePurchasedArchiveById(c, id, info); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
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

// @Summary Search materials
// @Security ApiKeyAuth
// @Tags materials
// @Description Поиск товара по наименованию
// @ID search-material
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, warehouse_id, item_id, name, article, product_category, total_quantity, volume, price_without_vat, total_without_vat, supplier_id, location, status, received_date, last_updated, min_stock_level, expiration_date, storage_cost, warehouse_section, incoming_delivery_number) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Param name query string true "name query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/search [GET]
func (h *Handler) searchMaterial(c *gin.Context) {
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

	query, err := parseNameQueryParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	mtrls, count, err := h.services.Materials.MaterialSearch(c.Request.Context(), domain.MaterialParams{
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
		SortField: field,
		Query:     query,
		CompanyId: info.CompanyId,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       mtrls,
		TotalCount: count,
	})
}

// @Summary Create material category
// @Security ApiKeyAuth
// @Tags materials category
// @Description Создание категории материала
// @ID create-material-category
// @Accept json
// @Produce json
// @Param input body domain.CreateMaterialCategory true "Необходимо указать данные категории материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/category [POST]
func (h *Handler) createCategory(c *gin.Context) {
	var inp domain.MaterialCategory
	if err := c.ShouldBindJSON(&inp); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.services.Category.Create(c.Request.Context(), domain.MaterialCategory{
		Name:        inp.Name,
		CompanyID:   info.CompanyId,
		Description: inp.Description,
		Slug:        inp.Slug,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		IsActive:    true,
		ImgURL:      inp.ImgURL,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newCreateSuccessIdResponse(c, id)
}

// @Summary Get material category
// @Security ApiKeyAuth
// @Tags materials category
// @Description Получение категории материала
// @ID get-material-category
// @Accept json
// @Produce json
// @Param id path int true "ID категории материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/category/{id} [GET]
func (h *Handler) getCategoryById(c *gin.Context) {
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

	category, err := h.services.Category.GetById(c.Request.Context(), id, info.CompanyId)
	if err != nil {
		if errors.Is(err, domain.ErrMaterialCategoryNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       category,
		TotalCount: 1,
	})
}

// @Summary Update material category
// @Security ApiKeyAuth
// @Tags materials category
// @Description Обновление категории материала
// @ID update-material-category
// @Accept json
// @Produce json
// @Param id path int true "ID категории материала"
// @Param input body domain.UpdateMaterialCategory true "Необходимо указать данные категории материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/category/{id} [PUT]
func (h *Handler) updateCategory(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var inp domain.UpdateMaterialCategory
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

	if err = h.services.Category.Update(c.Request.Context(), inp); err != nil {
		if errors.Is(err, domain.ErrMaterialCategoryNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Delete material category
// @Security ApiKeyAuth
// @Tags materials category
// @Description Удаление категории материала
// @ID delete-material-category
// @Accept json
// @Produce json
// @Param id path int true "ID категории материала"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/category/{id} [DELETE]
func (h *Handler) deleteCategory(c *gin.Context) {
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

	if err = h.services.Category.Delete(c, id, info.CompanyId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Get material category list
// @Security ApiKeyAuth
// @Tags materials category
// @Description Список категорий материалов
// @ID get-material-category-list
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, name, slug, created_at, updated_at, is_active) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/category [GET]
func (h *Handler) getCategoryList(c *gin.Context) {
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

	categories, count, err := h.services.Category.List(c.Request.Context(), domain.MaterialParams{
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
		Data:       categories,
		TotalCount: count,
	})
}

// @Summary Search material category
// @Security ApiKeyAuth
// @Tags materials category
// @Description Поиск категорий материалов
// @ID search-material-category
// @Accept json
// @Produce json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, name, slug, created_at, updated_at, is_active) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Param name query string true "offset query param"
// @Success 200 {object} []domain.MaterialCategory
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /materials/category/search [GET]
func (h *Handler) searchCategory(c *gin.Context) {
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

	name, err := parseNameQueryParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	categories, count, err := h.services.Category.Search(c.Request.Context(), domain.MaterialParams{
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
		SortField: field,
		CompanyId: info.CompanyId,
		Query:     name,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       categories,
		TotalCount: count,
	})
}
