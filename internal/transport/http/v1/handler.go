package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/service"
	"github.com/rusystem/crm-api/pkg/auth"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/rusystem/crm-api/tools"
	"strconv"
	"strings"
)

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
	cfg          *config.Config
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager, cfg *config.Config) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
		cfg:          cfg,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initAuthRoutes(v1)

		// warehouse routes
		h.initSupplierRoutes(v1)
		h.initWarehouseRoutes(v1)
		h.initMaterialsRoutes(v1)

		// accounts routes
		h.initCompanyRoutes(v1)
		h.initUserRoutes(v1)
		h.initSectionRoutes(v1)

		// geo route
		h.initGeoRoutes(v1)
	}
}

func parseOffsetQueryParam(c *gin.Context) (int64, error) {
	var offset int
	var err error

	offsetParam := c.Query("offset")
	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			return 0, domain.ErrInvalidOffsetParam
		}
	}

	if offset < 0 {
		return 0, domain.ErrInvalidOffsetParam
	}

	return int64(offset), nil
}

func parseLimitQueryParam(c *gin.Context) (int64, error) {
	var limit = 100
	var err error

	limitParam := c.Query("limit")
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return 0, domain.ErrInvalidLimitParam
		}
	}

	if limit <= 0 {
		return 0, domain.ErrInvalidLimitParam
	}

	return int64(limit), nil
}

var allowedFields = map[string]bool{
	"id":                       true,
	"username":                 true,
	"name":                     true,
	"address":                  true,
	"responsible_person":       true,
	"phone":                    true,
	"email":                    true,
	"max_capacity":             true,
	"current_occupancy":        true,
	"country":                  true,
	"region":                   true,
	"created_at":               true,
	"updated_at":               true,
	"position":                 true,
	"warehouse_id":             true,
	"item_id":                  true,
	"article":                  true,
	"product_category":         true,
	"total_quantity":           true,
	"volume":                   true,
	"price_without_vat":        true,
	"total_without_vat":        true,
	"supplier_id":              true,
	"location":                 true,
	"status":                   true,
	"received_date":            true,
	"last_updated":             true,
	"min_stock_level":          true,
	"expiration_date":          true,
	"storage_cost":             true,
	"warehouse_section":        true,
	"incoming_delivery_number": true,
	"slug":                     true,
	"is_active":                true,
	"name_ru":                  true,
	"name_en":                  true,
	"website":                  true,
	"is_approved":              true,
	"legal_address":            true,
	"actual_address":           true,
	"warehouse_address":        true,
	"contact_person":           true,
	"contract_number":          true,
	"product_categories":       true,
	"purchase_amount":          true,
	"balance":                  true,
	"product_types":            true,
	"tax_id":                   true,
	"registration_date":        true,
	"last_login":               true,
	"role":                     true,
}

func parseSortParam(c *gin.Context) (string, string, error) {
	sortParam := c.Query("sort")
	if sortParam == "" {
		sortParam = "asc"
	}

	if sortParam != "asc" && sortParam != "desc" {
		return "", "", domain.ErrInvalidSortParam
	}

	sortField := c.Query("sort_field")
	if sortField == "" {
		return "", "", domain.ErrInvalidSortFieldParam
	}

	// Проверяем, что поле сортировки разрешено
	if _, ok := allowedFields[sortField]; !ok {
		return "", "", domain.ErrInvalidSortFieldParam
	}

	return strings.ToUpper(sortParam), sortField, nil
}

func parseIdIntPathParam(c *gin.Context) (int64, error) {
	idParam := c.Param("id")
	if idParam == "" {
		return 0, domain.ErrInvalidIdParam
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, domain.ErrInvalidIdParam
	}

	if id <= 0 {
		return 0, domain.ErrInvalidIdParam
	}

	return int64(id), nil
}

func parseNameQueryParam(c *gin.Context) (string, error) {
	queryParam := c.Query("name")
	if queryParam == "" {
		return "", domain.ErrInvalidQueryParam
	}

	return queryParam, nil
}

func parseIdStringPathParam(c *gin.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", domain.ErrInvalidIdParam
	}

	return id, nil
}

func parseEmailPathParam(c *gin.Context) (string, error) {
	id := c.Param("email")
	if id == "" {
		return "", domain.ErrInvalidEmailParam
	}

	return id, nil
}

func parseCountryCodeStringPathParam(c *gin.Context) (string, error) {
	code := c.Query("country_code")
	if code == "" {
		return "", domain.ErrInvalidCountryCodeParam
	}

	code = strings.ToUpper(code)

	if !tools.IsValidCountryCode(code) {
		return "", domain.ErrInvalidCountryCodeParam
	}

	return code, nil
}

func parseAdminCodeStringPathParam(c *gin.Context) (string, error) {
	code := c.Query("admin_code")
	if code == "" {
		return "", domain.ErrInvalidAdminCodeParam
	}

	if !tools.IsValidAdminCode(code) {
		return "", domain.ErrInvalidAdminCodeParam
	}

	return code, nil
}
