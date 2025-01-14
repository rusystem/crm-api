package service

import (
	"github.com/rusystem/cache"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	"github.com/rusystem/crm-api/pkg/auth"
	"github.com/rusystem/crm-api/pkg/client/geonames"
)

type Config struct {
	Config       *config.Config
	Repo         *repository.Repository
	TokenManager auth.TokenManager
}

type Service struct {
	Auth          Auth
	Supplier      Supplier
	Warehouse     Warehouse
	User          User
	Company       Company
	Sections      Sections
	Materials     Materials
	Category      Category
	Geo           Geo
	UnitOfMeasure UnitOfMeasure
}

func New(cfg Config, gc *geonames.Client, cache *cache.MemoryCache) *Service {
	return &Service{
		Auth:          NewAuthServices(cfg.Config, cfg.Repo, cfg.TokenManager),
		Supplier:      NewSupplierService(cfg.Config, cfg.Repo),
		Warehouse:     NewWarehouseServices(cfg.Config, cfg.Repo),
		User:          NewUserServices(cfg.Config, cfg.Repo),
		Company:       NewCompanyService(cfg.Config, cfg.Repo),
		Sections:      NewSectionsService(cfg.Config, cfg.Repo),
		Materials:     NewMaterialsService(cfg.Config, cfg.Repo),
		Category:      NewMaterialCategoriesService(cfg.Config, cfg.Repo),
		Geo:           NewGeoService(cfg.Config, gc, cache),
		UnitOfMeasure: NewUnitOfMeasureService(cfg.Config, cfg.Repo),
	}
}
