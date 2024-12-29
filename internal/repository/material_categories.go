package repository

import (
	"context"
	"database/sql"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository/database"
	"github.com/rusystem/crm-api/pkg/domain"
)

type MaterialCategory interface {
	Create(ctx context.Context, category domain.MaterialCategory) (int64, error)
	GetById(ctx context.Context, id, companyId int64) (domain.MaterialCategory, error)
	Update(ctx context.Context, category domain.MaterialCategory) error
	Delete(ctx context.Context, id, companyId int64) error
	List(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error)
	Search(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error)
}

type MaterialCategoriesRepository struct {
	cfg *config.Config
	db  database.MaterialCategory
}

func NewMaterialCategoriesRepository(cfg *config.Config, db *sql.DB) *MaterialCategoriesRepository {
	return &MaterialCategoriesRepository{
		cfg: cfg,
		db:  database.NewMaterialCategoriesPostgresRepository(db),
	}
}

func (mc *MaterialCategoriesRepository) Create(ctx context.Context, category domain.MaterialCategory) (int64, error) {
	return mc.db.Create(ctx, category)
}

func (mc *MaterialCategoriesRepository) GetById(ctx context.Context, id, companyId int64) (domain.MaterialCategory, error) {
	return mc.db.GetById(ctx, id, companyId)
}

func (mc *MaterialCategoriesRepository) Update(ctx context.Context, category domain.MaterialCategory) error {
	return mc.db.Update(ctx, category)
}

func (mc *MaterialCategoriesRepository) Delete(ctx context.Context, id, companyId int64) error {
	return mc.db.Delete(ctx, id, companyId)
}

func (mc *MaterialCategoriesRepository) List(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error) {
	return mc.db.List(ctx, param)
}

func (mc *MaterialCategoriesRepository) Search(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error) {
	return mc.db.Search(ctx, param)
}
