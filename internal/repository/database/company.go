package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rusystem/crm-api/pkg/domain"
)

type Company interface {
	GetById(ctx context.Context, id int64) (domain.Company, error)
	IsExist(ctx context.Context, id int64) (bool, error)
	Create(ctx context.Context, company domain.Company) (int64, error)
	Update(ctx context.Context, company domain.Company) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, param domain.Param) ([]domain.Company, int64, error)
}

type CompanyDatabaseRepository struct {
	db *sql.DB
}

func NewCompanyDatabase(db *sql.DB) *CompanyDatabaseRepository {
	return &CompanyDatabaseRepository{
		db: db,
	}
}

func (cdr *CompanyDatabaseRepository) GetById(ctx context.Context, id int64) (domain.Company, error) {
	query := fmt.Sprintf(`SELECT id, name_ru, name_en, country, address, phone, email, website, is_active, created_at, updated_at, is_approved, timezone FROM %s WHERE id = $1`,
		domain.CompaniesTable)

	var company domain.Company
	err := cdr.db.QueryRowContext(ctx, query, id).Scan(
		&company.ID,
		&company.NameRu,
		&company.NameEn,
		&company.Country,
		&company.Address,
		&company.Phone,
		&company.Email,
		&company.Website,
		&company.IsActive,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.IsApproved,
		&company.Timezone,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return company, domain.ErrCompanyNotFound
		}

		return company, err
	}

	return company, nil
}

func (cdr *CompanyDatabaseRepository) IsExist(ctx context.Context, id int64) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)`, domain.CompaniesTable)

	var exists bool
	if err := cdr.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.ErrCompanyNotFound
		}

		return false, err
	}

	return exists, nil
}

func (cdr *CompanyDatabaseRepository) Create(ctx context.Context, company domain.Company) (int64, error) {
	var id int64
	query := fmt.Sprintf(`
		INSERT INTO %s
		(name_ru, name_en, country, address, phone, email, website, is_active, created_at, updated_at, is_approved, timezone)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id;
	`, domain.CompaniesTable)

	err := cdr.db.QueryRowContext(ctx, query,
		company.NameRu, company.NameEn, company.Country, company.Address, company.Phone, company.Email,
		company.Website, company.IsActive, company.CreatedAt, company.UpdatedAt, company.IsApproved, company.Timezone,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (cdr *CompanyDatabaseRepository) Update(ctx context.Context, company domain.Company) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET
		    name_ru = $1, name_en = $2, country = $3, address = $4, phone = $5, email = $6,
		    website = $7, is_active = $8, updated_at = $9, is_approved = $10, timezone = $11
		WHERE id = $12;
	`, domain.CompaniesTable)

	_, err := cdr.db.ExecContext(ctx, query,
		company.NameRu, company.NameEn, company.Country, company.Address, company.Phone, company.Email,
		company.Website, company.IsActive, company.UpdatedAt, company.IsApproved, company.Timezone,
		company.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (cdr *CompanyDatabaseRepository) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, domain.CompaniesTable)

	_, err := cdr.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (cdr *CompanyDatabaseRepository) List(ctx context.Context, param domain.Param) ([]domain.Company, int64, error) {
	var companies []domain.Company
	var count int64

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, domain.CompaniesTable)
	err := cdr.db.QueryRowContext(ctx, countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT 
		    id, name_ru, name_en, country, address, phone, email, website, 
		    is_active, created_at, updated_at, is_approved, timezone 
		FROM %s ORDER BY %s %s
		LIMIT $1 OFFSET $2;
	`, domain.CompaniesTable, param.SortField, param.Sort)

	rows, err := cdr.db.QueryContext(ctx, query, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var company domain.Company

		if err := rows.Scan(
			&company.ID, &company.NameRu, &company.NameEn, &company.Country, &company.Address, &company.Phone, &company.Email,
			&company.Website, &company.IsActive, &company.CreatedAt, &company.UpdatedAt, &company.IsApproved, &company.Timezone,
		); err != nil {
			return nil, 0, err
		}

		companies = append(companies, company)
	}

	return companies, count, nil
}
