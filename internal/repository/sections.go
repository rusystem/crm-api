package repository

import (
	"context"
	"database/sql"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository/database"
	"github.com/rusystem/crm-api/pkg/domain"
)

type Sections interface {
	GetById(ctx context.Context, id int64) (domain.Section, error)
	Create(ctx context.Context, section domain.Section) (int64, error)
	Update(ctx context.Context, section domain.Section) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, param domain.Param) ([]domain.Section, int64, error)
}

type SectionRepository struct {
	cfg  *config.Config
	psql database.Sections
}

func NewSectionsRepository(cfg *config.Config, db *sql.DB) *SectionRepository {
	return &SectionRepository{
		cfg:  cfg,
		psql: database.NewSectionsPostgresRepository(db),
	}
}

func (r *SectionRepository) GetById(ctx context.Context, id int64) (domain.Section, error) {
	return r.psql.GetById(ctx, id)
}

func (r *SectionRepository) Create(ctx context.Context, section domain.Section) (int64, error) {
	return r.psql.Create(ctx, section)
}

func (r *SectionRepository) Update(ctx context.Context, section domain.Section) error {
	return r.psql.Update(ctx, section)
}

func (r *SectionRepository) Delete(ctx context.Context, id int64) error {
	return r.psql.Delete(ctx, id)
}

func (r *SectionRepository) List(ctx context.Context, param domain.Param) ([]domain.Section, int64, error) {
	return r.psql.List(ctx, param)
}
