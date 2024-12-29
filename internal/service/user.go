package service

import (
	"context"
	"errors"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/rusystem/crm-api/tools"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	GetById(ctx context.Context, id int64, info domain.JWTInfo) (domain.User, error)
	UpdateProfile(ctx context.Context, user domain.UserProfileUpdate, info domain.JWTInfo) error
	Create(ctx context.Context, user domain.User) (int64, error)
	Update(ctx context.Context, user domain.UserUpdate, info domain.JWTInfo) error
	Delete(ctx context.Context, id int64, info domain.JWTInfo) error
	GetListByCompanyId(ctx context.Context, companyId int64, param domain.Param) ([]domain.UserResponse, int64, error)
}

type UserService struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewUserServices(cfg *config.Config, repo *repository.Repository) *UserService {
	return &UserService{
		cfg:  cfg,
		repo: repo,
	}
}

func (su *UserService) GetById(ctx context.Context, id int64, info domain.JWTInfo) (domain.User, error) {
	user, err := su.repo.User.GetById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	if user.CompanyID != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.User{}, domain.ErrNotAllowed
	}

	return user, nil
}

func (su *UserService) UpdateProfile(ctx context.Context, req domain.UserProfileUpdate, info domain.JWTInfo) error {
	user, err := su.repo.User.GetById(ctx, info.UserId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}

		return err
	}

	if user.CompanyID != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.ErrNotAllowed
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Phone != nil {
		user.Phone = *req.Phone
	}

	if req.Password != nil {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return domain.ErrUpdateUser
		}

		user.PasswordHash = string(hashedPass)
	}

	if req.Country != nil {
		user.Country = *req.Country
	}

	return su.repo.User.Update(ctx, user)
}

func (su *UserService) Create(ctx context.Context, user domain.User) (int64, error) {
	return su.repo.User.Create(ctx, user)
}

func (su *UserService) Update(ctx context.Context, req domain.UserUpdate, info domain.JWTInfo) error {
	user, err := su.repo.User.GetById(ctx, *req.ID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	if user.CompanyID != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.ErrNotAllowed
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Phone != nil {
		user.Phone = *req.Phone
	}

	if req.Password != nil {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return domain.ErrUpdateUser
		}

		user.PasswordHash = string(hashedPass)
	}

	if req.Language != nil {
		user.Language = *req.Language
	}

	if req.Country != nil {
		user.Country = *req.Country
	}

	if req.Position != nil {
		user.Position = *req.Position
	}

	if req.IsSendSystemNotification != nil {
		user.IsSendSystemNotification = *req.IsSendSystemNotification
	}

	if req.Sections != nil {
		if tools.IsFullAccessSection(*req.Sections) {
			return domain.ErrNotAllowed
		}

		user.Sections = *req.Sections
	}

	if req.Role != nil {
		if !tools.IsFullAccessSection(info.Sections) {
			return domain.ErrNotAllowed
		}

		user.Role = *req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if req.IsApproved != nil {
		user.IsApproved = *req.IsApproved
	}

	return su.repo.User.Update(ctx, user)
}

func (su *UserService) Delete(ctx context.Context, id int64, info domain.JWTInfo) error {
	user, err := su.repo.User.GetById(ctx, id)
	if err != nil {
		return err
	}

	if user.CompanyID != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.ErrNotAllowed
	}

	return su.repo.User.Delete(ctx, id)
}

func (su *UserService) GetListByCompanyId(ctx context.Context, companyId int64, param domain.Param) ([]domain.UserResponse, int64, error) {
	var resp []domain.UserResponse

	users, count, err := su.repo.User.GetListByCompanyId(ctx, companyId, param)
	if err != nil {
		return resp, 0, err
	}

	for _, v := range users {
		resp = append(resp, domain.UserResponse{
			ID:                       v.ID,
			CompanyID:                v.CompanyID,
			Username:                 v.Username,
			Name:                     v.Name,
			Email:                    v.Email,
			Phone:                    v.Phone,
			CreatedAt:                v.CreatedAt,
			UpdatedAt:                v.UpdatedAt,
			LastLogin:                v.LastLogin.Time,
			IsActive:                 v.IsActive,
			Role:                     v.Role,
			Language:                 v.Language,
			Country:                  v.Country,
			IsApproved:               v.IsApproved,
			IsSendSystemNotification: v.IsSendSystemNotification,
			Sections:                 v.Sections,
			Position:                 v.Position,
		})
	}

	return resp, count, nil
}
