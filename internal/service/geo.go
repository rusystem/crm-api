package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/rusystem/cache"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/pkg/client/geonames"
	"github.com/rusystem/crm-api/pkg/domain"
)

type Geo interface {
	CountryList(ctx context.Context, lang string) ([]domain.Country, error)
	RegionList(ctx context.Context, countryCode, lang string) ([]domain.Region, error)
	CityList(ctx context.Context, countryCode, adminCode, lang string) ([]domain.City, error)
}

type GeoService struct {
	cfg       *config.Config
	geoClient *geonames.Client
	cache     *cache.MemoryCache
}

func NewGeoService(cfg *config.Config, geoClient *geonames.Client, cache *cache.MemoryCache) *GeoService {
	return &GeoService{
		cfg:       cfg,
		geoClient: geoClient,
		cache:     cache,
	}
}

func (g *GeoService) CountryList(ctx context.Context, lang string) ([]domain.Country, error) {
	key := fmt.Sprintf("CountryList:%s", lang)

	cachedCountries, err := g.cache.Get(key)
	if err == nil {
		countries, ok := cachedCountries.([]domain.Country)
		if !ok {
			return nil, errors.New("can`t to cast country types")
		}

		return countries, nil
	}

	countries, err := g.geoClient.FetchCountries(ctx, lang)
	if err != nil {
		return nil, err
	}

	if err = g.cache.Set(key, countries, 0); err != nil {
		return nil, err
	}

	return countries, nil
}

func (g *GeoService) RegionList(ctx context.Context, countryCode, lang string) ([]domain.Region, error) {
	key := fmt.Sprintf("RegionList:%s%s", countryCode, lang)

	cachedRegions, err := g.cache.Get(key)
	if err == nil {
		regions, ok := cachedRegions.([]domain.Region)
		if !ok {
			return nil, errors.New("can`t to cast city types")
		}

		return regions, nil
	}

	regions, err := g.geoClient.FetchRegions(ctx, countryCode, lang)
	if err != nil {
		return nil, err
	}

	if err = g.cache.Set(key, regions, 0); err != nil {
		return nil, err
	}

	return regions, nil
}

func (g *GeoService) CityList(ctx context.Context, countryCode, adminCode, lang string) ([]domain.City, error) {
	key := fmt.Sprintf("CityList:%s%s%s", countryCode, adminCode, lang)

	cachedCities, err := g.cache.Get(key)
	if err == nil {
		cities, ok := cachedCities.([]domain.City)
		if !ok {
			return nil, errors.New("can`t to cast city types")
		}

		return cities, nil
	}

	cities, err := g.geoClient.FetchCitiesByRegion(ctx, countryCode, adminCode, lang)
	if err != nil {
		return nil, err
	}

	if err = g.cache.Set(key, cities, 0); err != nil {
		return nil, err
	}

	return cities, nil
}
