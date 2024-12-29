package repository

import (
	"context"
	"database/sql"
	"github.com/rusystem/cache"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository/database"
	"github.com/rusystem/crm-api/pkg/domain"
)

type Auth interface {
	CreateToken(ctx context.Context, token domain.RefreshSession) error
	DeleteToken(ctx context.Context, userId int64, token string) error
	DeleteUserTokens(ctx context.Context, userId, companyId int64) error
	GetSessionToken(ctx context.Context, refreshToken string) (domain.RefreshSession, error)
	GetUserWhiteIp(ctx context.Context, userId int64) ([]string, error)
}

type AuthRepository struct {
	cfg   *config.Config
	cache *cache.MemoryCache
	db    database.Auth
}

func NewAuthRepository(cfg *config.Config, cache *cache.MemoryCache, db *sql.DB) *AuthRepository {
	return &AuthRepository{
		cfg:   cfg,
		cache: cache,
		db:    database.NewAuthRepository(db),
	}
}

func (ar *AuthRepository) CreateToken(ctx context.Context, token domain.RefreshSession) error {
	return ar.db.CreateToken(ctx, token)
}

func (ar *AuthRepository) DeleteToken(ctx context.Context, userId int64, token string) error {
	return ar.db.DeleteToken(ctx, userId, token)
}

func (ar *AuthRepository) DeleteUserTokens(ctx context.Context, userId, companyId int64) error {
	return ar.db.DeleteUserTokens(ctx, userId, companyId)
}

func (ar *AuthRepository) GetSessionToken(ctx context.Context, refreshToken string) (domain.RefreshSession, error) {
	return ar.db.GetSessionToken(ctx, refreshToken)
}

func (ar *AuthRepository) GetUserWhiteIp(ctx context.Context, userId int64) ([]string, error) {
	return ar.db.GetUserWhiteIp(ctx, userId)
}
