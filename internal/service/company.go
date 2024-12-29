package service

import (
	"context"
	"errors"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	"github.com/rusystem/crm-api/pkg/domain"
	tools "github.com/rusystem/crm-api/tools"
	"time"
)

type Company interface {
	GetById(ctx context.Context, id int64) (domain.Company, error)
	Create(ctx context.Context, company domain.Company) (int64, error)
	Update(ctx context.Context, company domain.CompanyUpdate, info domain.JWTInfo) error
	Delete(ctx context.Context, id int64) error
	IsExist(ctx context.Context, id int64) (bool, error)
	List(ctx context.Context, param domain.Param) ([]domain.Company, int64, error)
}

type CompanyService struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewCompanyService(cfg *config.Config, repo *repository.Repository) *CompanyService {
	return &CompanyService{
		cfg:  cfg,
		repo: repo,
	}
}

func (c *CompanyService) GetById(ctx context.Context, id int64) (domain.Company, error) {
	return c.repo.Company.GetById(ctx, id)
}

func (c *CompanyService) Create(ctx context.Context, company domain.Company) (int64, error) {
	return c.repo.Company.Create(ctx, company)
}

func (c *CompanyService) Update(ctx context.Context, req domain.CompanyUpdate, info domain.JWTInfo) error {
	company, err := c.repo.Company.GetById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			return domain.ErrCompanyNotFound
		}

		return err
	}

	if company.ID != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.ErrNotAllowed
	}

	if req.NameRu != nil {
		company.NameRu = *req.NameRu
	}

	if req.NameEn != nil {
		company.NameEn = *req.NameEn
	}

	if req.Country != nil {
		company.Country = *req.Country
	}

	if req.Address != nil {
		company.Address = *req.Address
	}

	if req.Phone != nil {
		company.Phone = *req.Phone
	}

	if req.Email != nil {
		company.Email = *req.Email
	}

	if req.Website != nil {
		company.Website = *req.Website
	}

	if req.Timezone != nil {
		company.Timezone = *req.Timezone

		_, err := time.LoadLocation(*req.Timezone)
		if err != nil {
			return domain.ErrInvalidTimezone
		}
	}

	if (req.IsApproved != nil || req.IsActive != nil) && !tools.IsFullAccessSection(info.Sections) {
		return domain.ErrNotAllowed
	}

	if req.IsApproved != nil {
		company.IsApproved = *req.IsApproved
	}

	if req.IsActive != nil {
		company.IsActive = *req.IsActive
	}

	return c.repo.Company.Update(ctx, company)
}

func (c *CompanyService) Delete(ctx context.Context, id int64) error {
	return c.repo.Company.Delete(ctx, id)
}

func (c *CompanyService) IsExist(ctx context.Context, id int64) (bool, error) {
	return c.repo.Company.IsExist(ctx, id)
}

func (c *CompanyService) List(ctx context.Context, param domain.Param) ([]domain.Company, int64, error) {
	return c.repo.Company.List(ctx, param)
}
