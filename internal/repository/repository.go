package repository

import (
	"database/sql"
	"github.com/rusystem/cache"
	"github.com/rusystem/crm-api/internal/config"
)

type Repository struct {
	Auth             Auth
	User             User
	Company          Company
	MaterialCategory MaterialCategory
	Materials        Materials
	Sections         Sections
	Suppliers        Suppliers
	Warehouse        Warehouse
	UnitOfMeasure    UnitOfMeasure
}

func New(cfg *config.Config, cache *cache.MemoryCache, pc *sql.DB) *Repository {
	return &Repository{
		Auth:             NewAuthRepository(cfg, cache, pc),
		User:             NewUserRepository(cfg, cache, pc),
		Company:          NewCompanyRepository(cfg, cache, pc),
		MaterialCategory: NewMaterialCategoriesRepository(cfg, pc),
		Materials:        NewMaterialsRepository(cfg, pc),
		Sections:         NewSectionsRepository(cfg, pc),
		Suppliers:        NewSuppliersRepository(cfg, pc),
		Warehouse:        NewWarehouseRepository(cfg, pc),
		UnitOfMeasure:    NewUnitOfMeasureRepository(cfg, pc),
	}
}
