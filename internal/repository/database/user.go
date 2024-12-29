package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/rusystem/crm-api/pkg/domain"
	"time"
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

type UserDatabaseRepository struct {
	db *sql.DB
}

func NewUserDatabase(db *sql.DB) *UserDatabaseRepository {
	return &UserDatabaseRepository{
		db: db,
	}
}

func (udr *UserDatabaseRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	query := fmt.Sprintf(`
        SELECT id, company_id, username, email, phone, password_hash, created_at, updated_at, last_login, is_active,
               role, language, country, is_approved, is_send_system_notification, sections, position
        FROM %s
        WHERE username = $1`, domain.UsersTable)

	var user domain.User
	var b []byte
	err := udr.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.CompanyID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.IsActive,
		&user.Role,
		&user.Language,
		&user.Country,
		&user.IsApproved,
		&user.IsSendSystemNotification,
		&b,
		&user.Position,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	if err = json.Unmarshal(b, &user.Sections); err != nil {
		return domain.User{}, fmt.Errorf("error unmarshalling sections: %v", err)
	}

	return user, nil
}

func (udr *UserDatabaseRepository) GetSections(ctx context.Context, id int64) ([]string, error) {
	sections := make([]string, 0)

	query := fmt.Sprintf(`
        SELECT sections
        FROM %s
        WHERE id = $1`, domain.UsersTable)

	var b []byte
	err := udr.db.QueryRowContext(ctx, query, id).Scan(&b)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sections, domain.ErrUserNotFound
		}

		return sections, err
	}

	if err = json.Unmarshal(b, &sections); err != nil {
		return sections, fmt.Errorf("error unmarshalling sections: %v", err)
	}

	return sections, nil
}

func (udr *UserDatabaseRepository) GetById(ctx context.Context, id int64) (domain.User, error) {
	query := fmt.Sprintf(`
        SELECT id, company_id, username, email, phone, password_hash, created_at, updated_at, last_login, is_active,
               role, language, country, is_approved, is_send_system_notification, sections, position
        FROM %s
        WHERE id = $1`, domain.UsersTable)

	var user domain.User
	var b []byte
	err := udr.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CompanyID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.IsActive,
		&user.Role,
		&user.Language,
		&user.Country,
		&user.IsApproved,
		&user.IsSendSystemNotification,
		&b,
		&user.Position,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	if err = json.Unmarshal(b, &user.Sections); err != nil {
		return domain.User{}, fmt.Errorf("error unmarshalling sections: %v", err)
	}

	return user, nil
}

func (udr *UserDatabaseRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`UPDATE %s SET last_login = $1 WHERE id = $2`, domain.UsersTable)

	_, err := udr.db.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	return nil
}

func (udr *UserDatabaseRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	query := fmt.Sprintf(`
        INSERT INTO %s
        (company_id, username, name, email, phone, password_hash, created_at, updated_at, last_login, is_active,
         role, language, country, is_approved, is_send_system_notification, sections, position)
        VALUES
        ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
        RETURNING id`, domain.UsersTable)

	sectionsJSON, err := json.Marshal(user.Sections)
	if err != nil {
		return 0, err
	}

	var id int64
	err = udr.db.QueryRowContext(ctx, query,
		user.CompanyID,
		user.Username,
		user.Name,
		user.Email,
		user.Phone,
		user.PasswordHash,
		time.Now().UTC(),
		time.Now().UTC(),
		user.LastLogin,
		user.IsActive,
		user.Role,
		user.Language,
		user.Country,
		user.IsApproved,
		user.IsSendSystemNotification,
		sectionsJSON,
		user.Position,
	).Scan(&id)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return 0, domain.ErrUserAlreadyExists
			}
		}

		return 0, err
	}

	return id, nil
}

func (udr *UserDatabaseRepository) Update(ctx context.Context, user domain.User) error {
	sectionsJSON, err := json.Marshal(user.Sections)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET
		    company_id = $1, username = $2, name = $3, email = $4, phone = $5, password_hash = $6, updated_at = $7,
		    last_login = $8, is_active = $9, role = $10, language = $11, country = $12, is_approved = $13,
		    is_send_system_notification = $14, sections = $15, position = $16
		WHERE id = $17
		`, domain.UsersTable)

	_, err = udr.db.ExecContext(ctx, query,
		user.CompanyID,
		user.Username,
		user.Name,
		user.Email,
		user.Phone,
		user.PasswordHash,
		time.Now().UTC(),
		user.LastLogin,
		user.IsActive,
		user.Role,
		user.Language,
		user.Country,
		user.IsApproved,
		user.IsSendSystemNotification,
		sectionsJSON,
		user.Position,
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (udr *UserDatabaseRepository) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, domain.UsersTable)

	_, err := udr.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (udr *UserDatabaseRepository) GetListByCompanyId(ctx context.Context, companyId int64, param domain.Param) ([]domain.User, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE company_id = $1
	`, domain.UsersTable)

	err := udr.db.QueryRowContext(ctx, countQuery, companyId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT 
		    id, company_id, username, name, email, phone, password_hash, created_at, 
		    updated_at, last_login, is_active, role, language, country, 
		    is_approved, is_send_system_notification, sections, position
		FROM %s
		WHERE company_id = $1 ORDER BY %s %s
		LIMIT $2 OFFSET $3
		`, domain.UsersTable, param.SortField, param.Sort)

	var users []domain.User

	rows, err := udr.db.QueryContext(ctx, query, companyId, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var user domain.User
		var b []byte

		if err := rows.Scan(
			&user.ID, &user.CompanyID, &user.Username, &user.Name, &user.Email, &user.Phone, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
			&user.LastLogin, &user.IsActive, &user.Role, &user.Language, &user.Country, &user.IsApproved, &user.IsSendSystemNotification,
			&b, &user.Position,
		); err != nil {
			return nil, 0, err
		}

		if err = json.Unmarshal(b, &user.Sections); err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}
