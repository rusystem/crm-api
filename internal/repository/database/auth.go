package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rusystem/crm-api/pkg/domain"
)

type Auth interface {
	CreateToken(ctx context.Context, token domain.RefreshSession) error
	DeleteToken(ctx context.Context, userId int64, token string) error
	DeleteUserTokens(ctx context.Context, userId, companyId int64) error
	GetSessionToken(ctx context.Context, refreshToken string) (domain.RefreshSession, error)
	GetUserWhiteIp(ctx context.Context, userId int64) ([]string, error)
}

type AuthDatabaseRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthDatabaseRepository {
	return &AuthDatabaseRepository{
		db: db,
	}
}

func (ar *AuthDatabaseRepository) CreateToken(ctx context.Context, token domain.RefreshSession) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, company_id, roles, token, expires_at, ip) VALUES ($1, $2, $3, $4, $5, $6);",
		domain.RefreshTokensTable)

	_, err := ar.db.ExecContext(ctx, query, token.UserID, token.CompanyID, token.Role, token.Token, token.ExpiresAt, token.Ip)
	if err != nil {
		return fmt.Errorf("could not insert token data: %v", err)
	}

	return nil
}

func (ar *AuthDatabaseRepository) DeleteToken(ctx context.Context, userId int64, token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND token = $2;", domain.RefreshTokensTable)

	_, err := ar.db.ExecContext(ctx, query, userId, token)
	if err != nil {
		return fmt.Errorf("could not delete token data: %v", err)
	}

	return nil
}

func (ar *AuthDatabaseRepository) DeleteUserTokens(ctx context.Context, userId, companyId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND company_id = $2;", domain.RefreshTokensTable)

	_, err := ar.db.ExecContext(ctx, query, userId, companyId)
	if err != nil {
		return fmt.Errorf("could not delete token data: %v", err)
	}

	return nil
}

func (ar *AuthDatabaseRepository) GetSessionToken(ctx context.Context, refreshToken string) (domain.RefreshSession, error) {
	query := fmt.Sprintf("SELECT id, user_id, company_id, roles, token, expires_at, ip FROM %s WHERE token = $1",
		domain.RefreshTokensTable)

	row := ar.db.QueryRowContext(ctx, query, refreshToken)

	var session domain.RefreshSession

	if err := row.Scan(&session.ID, &session.UserID, &session.CompanyID, &session.Role, &session.Token, &session.ExpiresAt, &session.Ip); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.RefreshSession{}, domain.ErrRefreshTokenNotFound
		}

		return domain.RefreshSession{}, err
	}

	return session, nil
}

func (ar *AuthDatabaseRepository) GetUserWhiteIp(ctx context.Context, userId int64) ([]string, error) {
	var ips []string

	query := fmt.Sprintf("SELECT ip FROM %s WHERE user_id = $1", domain.RefreshTokensTable)

	rows, err := ar.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		ips = append(ips, ip)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return ips, nil
}
