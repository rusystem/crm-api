package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	"net/http"
)

func (h *Handler) initGeoRoutes(api *gin.RouterGroup) {
	geo := api.Group("/geo", h.userIdentity)
	{
		country := geo.Group("/country")
		{
			country.GET("/list", h.getCountryList)
		}

		region := geo.Group("/region")
		{
			region.GET("/list", h.getRegionList)
		}

		city := geo.Group("/city")
		{
			city.GET("/list", h.getCityList)
		}
	}
}

// @Summary Get country list
// @Security ApiKeyAuth
// @Tags geo
// @Description Список стран
// @ID get-country-list
// @Accept json
// @Produce json
// @Success 200 {object} domain.Country
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /geo/country/list [GET]
func (h *Handler) getCountryList(c *gin.Context) {
	countries, err := h.services.Geo.CountryList(c, "ru") // todo учесть при локализации
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       countries,
		TotalCount: int64(len(countries)),
	})
}

// @Summary Get region list
// @Security ApiKeyAuth
// @Tags geo
// @Description Список регионов/областей
// @ID get-region-list
// @Accept json
// @Produce json
// @Param country_id query string true "country id query param" example(RU)
// @Success 200 {object} domain.Region
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /geo/region/list [GET]
func (h *Handler) getRegionList(c *gin.Context) {
	countryId, err := parseCountryCodeStringPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	regions, err := h.services.Geo.RegionList(c, countryId, "ru") //todo учесть локализацию
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       regions,
		TotalCount: int64(len(regions)),
	})
}

// @Summary Get city list
// @Security ApiKeyAuth
// @Tags geo
// @Description Список городов
// @ID get-city-list
// @Accept json
// @Produce json
// @Param country_id query string true "country id query param" example(RU)
// @Param region_id query string true "region id query param" example(65)
// @Success 200 {object} domain.City
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /geo/city/list [GET]
func (h *Handler) getCityList(c *gin.Context) {
	countryId, err := parseCountryCodeStringPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	regionId, err := parseRegionStringPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	cities, err := h.services.Geo.CityList(c, countryId, regionId, "ru") //todo учесть локализацию
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       cities,
		TotalCount: int64(len(cities)),
	})
}
