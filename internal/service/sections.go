package service

import (
	"context"
	"errors"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/rusystem/crm-api/tools"
)

type Sections interface {
	GetById(ctx context.Context, id int64) (domain.Section, error)
	Create(ctx context.Context, section domain.Section) (int64, error)
	Update(ctx context.Context, section domain.Section) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, info domain.JWTInfo, param domain.Param) ([]domain.Section, int64, error)
}

type SectionsService struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewSectionsService(cfg *config.Config, repo *repository.Repository) *SectionsService {
	return &SectionsService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *SectionsService) GetById(ctx context.Context, id int64) (domain.Section, error) {
	return s.repo.Sections.GetById(ctx, id)
}

func (s *SectionsService) Create(ctx context.Context, section domain.Section) (int64, error) {
	return s.repo.Sections.Create(ctx, section)
}

func (s *SectionsService) Update(ctx context.Context, section domain.Section) error {
	oldSection, err := s.repo.Sections.GetById(ctx, section.Id)
	if err != nil {
		return err
	}

	if oldSection.Name == domain.SectionFullAllAccess {
		return domain.ErrNotAllowed
	}

	return s.repo.Sections.Update(ctx, section)
}

func (s *SectionsService) Delete(ctx context.Context, id int64) error {
	section, err := s.repo.Sections.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrSectionNotFound) {
			return domain.ErrSectionNotFound
		}

		return err
	}

	if section.Name == domain.SectionFullAllAccess {
		return domain.ErrNotAllowed
	}

	return s.repo.Sections.Delete(ctx, id)
}

func (s *SectionsService) List(ctx context.Context, info domain.JWTInfo, param domain.Param) ([]domain.Section, int64, error) {
	sections, count, err := s.repo.Sections.List(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	if !tools.IsFullAccessSection(info.Sections) {
		tools.RemoveFullAccessSection(sections, domain.SectionFullAllAccess)
	}

	return sections, count, nil
}
