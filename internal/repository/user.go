package repository

import (
	"context"
	"database/sql"
	"github.com/rusystem/cache"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository/database"
	"github.com/rusystem/crm-api/pkg/domain"
)

type User interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	GetSections(ctx context.Context, id int64) ([]string, error)
	GetById(ctx context.Context, id int64) (domain.User, error)
	UpdateLastLogin(ctx context.Context, id int64) error
	Create(ctx context.Context, user domain.User) (int64, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id int64) error
	GetListByCompanyId(ctx context.Context, companyId int64, param domain.Param) ([]domain.User, int64, error)
}

type UserRepository struct {
	cfg   *config.Config
	cache *cache.MemoryCache
	db    database.User
}

func NewUserRepository(cfg *config.Config, cache *cache.MemoryCache, db *sql.DB) *UserRepository {
	return &UserRepository{
		cfg:   cfg,
		cache: cache,
		db:    database.NewUserDatabase(db),
	}
}

func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return ur.db.GetByUsername(ctx, username)
}

func (ur *UserRepository) GetSections(ctx context.Context, id int64) ([]string, error) {
	return ur.db.GetSections(ctx, id)
}

func (ur *UserRepository) GetById(ctx context.Context, id int64) (domain.User, error) {
	return ur.db.GetById(ctx, id)
}

func (ur *UserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	return ur.db.UpdateLastLogin(ctx, id)
}

func (ur *UserRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	return ur.db.Create(ctx, user)
}

func (ur *UserRepository) Update(ctx context.Context, user domain.User) error {
	return ur.db.Update(ctx, user)
}

func (ur *UserRepository) Delete(ctx context.Context, id int64) error {
	return ur.db.Delete(ctx, id)
}

func (ur *UserRepository) GetListByCompanyId(ctx context.Context, companyId int64, param domain.Param) ([]domain.User, int64, error) {
	return ur.db.GetListByCompanyId(ctx, companyId, param)
}
