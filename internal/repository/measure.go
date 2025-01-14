package repository

import (
	"context"
	"database/sql"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository/database"
	"github.com/rusystem/crm-api/pkg/domain"
)

type UnitOfMeasure interface {
	Create(ctx context.Context, measure domain.UnitOfMeasure) (int64, error)
	Update(ctx context.Context, measure domain.UnitOfMeasure) error
	Delete(ctx context.Context, id, companyId int64) error
	GetById(ctx context.Context, id, companyId int64) (domain.UnitOfMeasure, error)
	List(ctx context.Context, param domain.Param) ([]domain.UnitOfMeasure, int64, error)
}

type UnitOfMeasureRepository struct {
	cfg *config.Config
	db  database.UnitOfMeasure
}

func NewUnitOfMeasureRepository(cfg *config.Config, db *sql.DB) *UnitOfMeasureRepository {
	return &UnitOfMeasureRepository{
		cfg: cfg,
		db:  database.NewUnitOfMeasurePostgresRepository(db),
	}
}

func (umr *UnitOfMeasureRepository) Create(ctx context.Context, measure domain.UnitOfMeasure) (int64, error) {
	return umr.db.Create(ctx, measure)
}

func (umr *UnitOfMeasureRepository) Update(ctx context.Context, measure domain.UnitOfMeasure) error {
	return umr.db.Update(ctx, measure)
}

func (umr *UnitOfMeasureRepository) Delete(ctx context.Context, id, companyId int64) error {
	return umr.db.Delete(ctx, id, companyId)
}

func (umr *UnitOfMeasureRepository) GetById(ctx context.Context, id, companyId int64) (domain.UnitOfMeasure, error) {
	return umr.db.GetById(ctx, id, companyId)
}

func (umr *UnitOfMeasureRepository) List(ctx context.Context, param domain.Param) ([]domain.UnitOfMeasure, int64, error) {
	return umr.db.List(ctx, param)
}
