package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rusystem/crm-api/pkg/domain"
)

type Sections interface {
	GetById(ctx context.Context, id int64) (domain.Section, error)
	Create(ctx context.Context, section domain.Section) (int64, error)
	Update(ctx context.Context, section domain.Section) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, param domain.Param) ([]domain.Section, int64, error)
}

type SectionPostgresRepository struct {
	db *sql.DB
}

func NewSectionsPostgresRepository(db *sql.DB) *SectionPostgresRepository {
	return &SectionPostgresRepository{db: db}
}

func (spr *SectionPostgresRepository) GetById(ctx context.Context, id int64) (domain.Section, error) {
	var section domain.Section

	query := fmt.Sprintf("SELECT id, name FROM %s WHERE id = $1", domain.SectionsTable)

	if err := spr.db.QueryRowContext(ctx, query, id).Scan(&section.Id, &section.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Section{}, domain.ErrSectionNotFound
		}

		return domain.Section{}, err
	}

	return section, nil
}

func (spr *SectionPostgresRepository) Create(ctx context.Context, section domain.Section) (int64, error) {
	var id int64

	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", domain.SectionsTable)

	if err := spr.db.QueryRowContext(ctx, query, section.Name).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (spr *SectionPostgresRepository) Update(ctx context.Context, section domain.Section) error {
	query := fmt.Sprintf("UPDATE %s SET name = $1 WHERE id = $2", domain.SectionsTable)

	_, err := spr.db.ExecContext(ctx, query, section.Name, section.Id)
	if err != nil {
		return err
	}

	return nil
}

func (spr *SectionPostgresRepository) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", domain.SectionsTable)

	_, err := spr.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (spr *SectionPostgresRepository) List(ctx context.Context, param domain.Param) ([]domain.Section, int64, error) {
	var sections []domain.Section
	var totalCount int64

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", domain.SectionsTable)
	err := spr.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT id, name FROM %s ORDER BY %s %s LIMIT $1 OFFSET $2",
		domain.SectionsTable, param.SortField, param.Sort)
	rows, err := spr.db.QueryContext(ctx, query, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var section domain.Section
		if err := rows.Scan(&section.Id, &section.Name); err != nil {
			return nil, 0, err
		}

		sections = append(sections, section)
	}

	return sections, totalCount, nil
}
