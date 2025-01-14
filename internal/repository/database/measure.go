package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rusystem/crm-api/pkg/domain"
)

type UnitOfMeasure interface {
	Create(ctx context.Context, measure domain.UnitOfMeasure) (int64, error)
	Update(ctx context.Context, measure domain.UnitOfMeasure) error
	Delete(ctx context.Context, id, companyId int64) error
	GetById(ctx context.Context, id, companyId int64) (domain.UnitOfMeasure, error)
	List(ctx context.Context, param domain.Param) ([]domain.UnitOfMeasure, int64, error)
}

type UnitOfMeasurePostgresRepository struct {
	psql *sql.DB
}

func NewUnitOfMeasurePostgresRepository(psql *sql.DB) *UnitOfMeasurePostgresRepository {
	return &UnitOfMeasurePostgresRepository{
		psql: psql,
	}
}

func (umr *UnitOfMeasurePostgresRepository) Create(ctx context.Context, m domain.UnitOfMeasure) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (name, name_en, abbreviation, description, company_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		domain.UnitsOfMeasureTable)

	var id int64
	if err := umr.psql.QueryRowContext(ctx, query,
		m.Name, m.NameEn, m.Abbreviation, m.Description, m.CompanyID,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert unit of measure: %v", err)
	}

	return id, nil
}

func (umr *UnitOfMeasurePostgresRepository) Update(ctx context.Context, m domain.UnitOfMeasure) error {
	query := fmt.Sprintf(`
		UPDATE %s 
		SET
			name = $1, name_en = $2, abbreviation = $3, description = $4
		WHERE id = $5 AND company_id = $6`,
		domain.UnitsOfMeasureTable)

	_, err := umr.psql.ExecContext(ctx, query,
		m.Name, m.NameEn, m.Abbreviation, m.Description, m.ID, m.CompanyID,
	)

	return err
}

func (umr *UnitOfMeasurePostgresRepository) Delete(ctx context.Context, id, companyId int64) error {
	_, err := umr.psql.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND company_id = $2",
		domain.UnitsOfMeasureTable), id, companyId)
	return err
}

func (umr *UnitOfMeasurePostgresRepository) GetById(ctx context.Context, id, companyId int64) (domain.UnitOfMeasure, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, name, name_en, abbreviation, description, company_id
		FROM %s WHERE id = $1 AND company_id = $2`,
		domain.UnitsOfMeasureTable)

	var m domain.UnitOfMeasure
	if err := umr.psql.QueryRowContext(ctx, query, id, companyId).Scan(
		&m.ID, &m.Name, &m.NameEn, &m.Abbreviation, &m.Description, &m.CompanyID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.UnitOfMeasure{}, domain.ErrUnitOfMeasureNotFound
		}

		return domain.UnitOfMeasure{}, err
	}

	return m, nil
}

func (umr *UnitOfMeasurePostgresRepository) List(ctx context.Context, param domain.Param) ([]domain.UnitOfMeasure, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE company_id = $1`,
		domain.UnitsOfMeasureTable)

	err := umr.psql.QueryRowContext(ctx, countQuery, param.CompanyId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT 
			id, name, name_en, abbreviation, description, company_id
		FROM %s WHERE company_id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3`,
		domain.UnitsOfMeasureTable, param.SortField, param.Sort)

	rows, err := umr.psql.QueryContext(ctx, query, param.CompanyId, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var measures []domain.UnitOfMeasure

	for rows.Next() {
		var m domain.UnitOfMeasure
		if err = rows.Scan(
			&m.ID, &m.Name, &m.NameEn, &m.Abbreviation, &m.Description, &m.CompanyID,
		); err != nil {
			return nil, 0, err
		}

		measures = append(measures, m)
	}

	return measures, totalCount, nil
}
