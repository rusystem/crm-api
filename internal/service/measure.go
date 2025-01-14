package service

import (
	"context"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	"github.com/rusystem/crm-api/pkg/domain"
)

type UnitOfMeasure interface {
	Create(ctx context.Context, measure domain.UnitOfMeasure) (int64, error)
	Update(ctx context.Context, measure domain.UpdateUnitOfMeasure) error
	Delete(ctx context.Context, id, companyId int64) error
	GetById(ctx context.Context, id, companyId int64) (domain.UnitOfMeasure, error)
	List(ctx context.Context, param domain.Param) ([]domain.UnitOfMeasure, int64, error)
}

type UnitOfMeasureService struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewUnitOfMeasureService(cfg *config.Config, repo *repository.Repository) *UnitOfMeasureService {
	return &UnitOfMeasureService{
		cfg:  cfg,
		repo: repo,
	}
}

func (ums *UnitOfMeasureService) Create(ctx context.Context, measure domain.UnitOfMeasure) (int64, error) {
	return ums.repo.UnitOfMeasure.Create(ctx, measure)
}

func (ums *UnitOfMeasureService) Update(ctx context.Context, inp domain.UpdateUnitOfMeasure) error {
	measure, err := ums.repo.UnitOfMeasure.GetById(ctx, inp.ID, inp.CompanyID)
	if err != nil {
		return err
	}

	if inp.Name != nil {
		measure.Name = *inp.Name
	}

	if inp.NameEn != nil {
		measure.NameEn = *inp.NameEn
	}

	if inp.Abbreviation != nil {
		measure.Abbreviation = *inp.Abbreviation
	}

	if inp.Description != nil {
		measure.Description = *inp.Description
	}

	return ums.repo.UnitOfMeasure.Update(ctx, measure)
}

func (ums *UnitOfMeasureService) Delete(ctx context.Context, id, companyId int64) error {
	return ums.repo.UnitOfMeasure.Delete(ctx, id, companyId)
}

func (ums *UnitOfMeasureService) GetById(ctx context.Context, id, companyId int64) (domain.UnitOfMeasure, error) {
	return ums.repo.UnitOfMeasure.GetById(ctx, id, companyId)
}

func (ums *UnitOfMeasureService) List(ctx context.Context, param domain.Param) ([]domain.UnitOfMeasure, int64, error) {
	return ums.repo.UnitOfMeasure.List(ctx, param)
}
