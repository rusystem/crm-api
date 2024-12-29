package service

import (
	"context"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	"github.com/rusystem/crm-api/pkg/domain"
	"time"
)

type Category interface {
	Create(ctx context.Context, category domain.MaterialCategory) (int64, error)
	GetById(ctx context.Context, id, companyId int64) (domain.MaterialCategory, error)
	Update(ctx context.Context, inp domain.UpdateMaterialCategory) error
	Delete(ctx context.Context, id, companyId int64) error
	List(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error)
	Search(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error)
}

type MaterialCategoriesService struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewMaterialCategoriesService(cfg *config.Config, repo *repository.Repository) *MaterialCategoriesService {
	return &MaterialCategoriesService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *MaterialCategoriesService) Create(ctx context.Context, category domain.MaterialCategory) (int64, error) {
	return s.repo.MaterialCategory.Create(ctx, category)
}

func (s *MaterialCategoriesService) GetById(ctx context.Context, id, companyId int64) (domain.MaterialCategory, error) {
	return s.repo.MaterialCategory.GetById(ctx, id, companyId)
}

func (s *MaterialCategoriesService) Update(ctx context.Context, inp domain.UpdateMaterialCategory) error {
	category, err := s.repo.MaterialCategory.GetById(ctx, inp.ID, inp.CompanyID)
	if err != nil {
		return err
	}

	if inp.Name != nil {
		category.Name = *inp.Name
	}

	if inp.Description != nil {
		category.Description = *inp.Description
	}

	if inp.Slug != nil {
		category.Slug = *inp.Slug
	}

	if inp.IsActive != nil {
		category.IsActive = *inp.IsActive
	}

	if inp.ImgURL != nil {
		category.ImgURL = *inp.ImgURL
	}

	category.UpdatedAt = time.Now().UTC()

	return s.repo.MaterialCategory.Update(ctx, category)
}

func (s *MaterialCategoriesService) Delete(ctx context.Context, id, companyId int64) error {
	return s.repo.MaterialCategory.Delete(ctx, id, companyId)
}

func (s *MaterialCategoriesService) List(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error) {
	return s.repo.MaterialCategory.List(ctx, param)
}

func (s *MaterialCategoriesService) Search(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error) {
	return s.repo.MaterialCategory.Search(ctx, param)
}
