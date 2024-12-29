package repository

import (
	"context"
	"database/sql"
	"github.com/rusystem/cache"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository/database"
	"github.com/rusystem/crm-api/pkg/domain"
)

type Company interface {
	GetById(ctx context.Context, id int64) (domain.Company, error)
	IsExist(ctx context.Context, id int64) (bool, error)
	Create(ctx context.Context, company domain.Company) (int64, error)
	Update(ctx context.Context, company domain.Company) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, param domain.Param) ([]domain.Company, int64, error)
}

type CompanyRepository struct {
	cfg   *config.Config
	cache *cache.MemoryCache
	db    database.Company
}

func NewCompanyRepository(cfg *config.Config, cache *cache.MemoryCache, db *sql.DB) *CompanyRepository {
	return &CompanyRepository{
		cfg:   cfg,
		cache: cache,
		db:    database.NewCompanyDatabase(db),
	}
}

func (c *CompanyRepository) GetById(ctx context.Context, id int64) (domain.Company, error) {
	return c.db.GetById(ctx, id)
}

func (c *CompanyRepository) IsExist(ctx context.Context, id int64) (bool, error) {
	return c.db.IsExist(ctx, id)
}

func (c *CompanyRepository) Create(ctx context.Context, company domain.Company) (int64, error) {
	return c.db.Create(ctx, company)
}

func (c *CompanyRepository) Update(ctx context.Context, company domain.Company) error {
	return c.db.Update(ctx, company)
}

func (c *CompanyRepository) Delete(ctx context.Context, id int64) error {
	return c.db.Delete(ctx, id)
}

func (c *CompanyRepository) List(ctx context.Context, param domain.Param) ([]domain.Company, int64, error) {
	return c.db.List(ctx, param)
}
