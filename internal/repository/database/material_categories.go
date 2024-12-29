package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rusystem/crm-api/pkg/domain"
)

type MaterialCategory interface {
	Create(ctx context.Context, category domain.MaterialCategory) (int64, error)
	GetById(ctx context.Context, id, companyId int64) (domain.MaterialCategory, error)
	Update(ctx context.Context, category domain.MaterialCategory) error
	Delete(ctx context.Context, id, companyId int64) error
	List(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error)
	Search(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error)
}

type MaterialCategoriesPostgresRepository struct {
	psql *sql.DB
}

func NewMaterialCategoriesPostgresRepository(psql *sql.DB) *MaterialCategoriesPostgresRepository {
	return &MaterialCategoriesPostgresRepository{
		psql: psql,
	}
}

func (mc *MaterialCategoriesPostgresRepository) Create(ctx context.Context, c domain.MaterialCategory) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (name, company_id, description, slug, created_at, updated_at, is_active, img_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		domain.TableMaterialCategories)

	var id int64
	if err := mc.psql.QueryRowContext(ctx, query,
		c.Name, c.CompanyID, c.Description, c.Slug, c.CreatedAt, c.UpdatedAt, c.IsActive, c.ImgURL,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert material category: %v", err)
	}

	return id, nil
}

func (mc *MaterialCategoriesPostgresRepository) GetById(ctx context.Context, id, companyId int64) (domain.MaterialCategory, error) {
	query := fmt.Sprintf(`
		SELECT 
		    id, name, company_id, description, slug, created_at, updated_at, is_active, img_url 
		FROM %s WHERE id = $1 AND company_id = $2`,
		domain.TableMaterialCategories)

	var c domain.MaterialCategory
	if err := mc.psql.QueryRowContext(ctx, query, id, companyId).Scan(
		&c.ID, &c.Name, &c.CompanyID, &c.Description, &c.Slug, &c.CreatedAt, &c.UpdatedAt, &c.IsActive, &c.ImgURL,
	); err != nil {
		return domain.MaterialCategory{}, err
	}

	return c, nil
}

func (mc *MaterialCategoriesPostgresRepository) Update(ctx context.Context, c domain.MaterialCategory) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET
			name = $1, description = $2, slug = $3, created_at = $4, updated_at = $5, is_active = $6, img_url = $7
		WHERE id = $8 AND company_id = $9`,
		domain.TableMaterialCategories)

	_, err := mc.psql.ExecContext(ctx, query,
		c.Name, c.Description, c.Slug, c.CreatedAt, c.UpdatedAt, c.IsActive, c.ImgURL, c.ID, c.CompanyID,
	)

	return err
}

func (mc *MaterialCategoriesPostgresRepository) Delete(ctx context.Context, id, companyId int64) error {
	_, err := mc.psql.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND company_id = $2",
		domain.TableMaterialCategories), id, companyId)
	return err
}

func (mc *MaterialCategoriesPostgresRepository) List(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE company_id = $1`,
		domain.TableMaterialCategories)

	err := mc.psql.QueryRowContext(ctx, countQuery, param.CompanyId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT 
		    id, name, company_id, description, slug, created_at, updated_at, is_active, img_url 
		FROM %s WHERE company_id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3`,
		domain.TableMaterialCategories, param.SortField, param.Sort)

	rows, err := mc.psql.QueryContext(ctx, query, param.CompanyId, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var categories []domain.MaterialCategory

	for rows.Next() {
		var c domain.MaterialCategory
		if err = rows.Scan(
			&c.ID, &c.Name, &c.CompanyID, &c.Description, &c.Slug, &c.CreatedAt, &c.UpdatedAt, &c.IsActive, &c.ImgURL,
		); err != nil {
			return nil, 0, err
		}

		categories = append(categories, c)
	}

	return categories, totalCount, nil
}

func (mc *MaterialCategoriesPostgresRepository) Search(ctx context.Context, param domain.MaterialParams) ([]domain.MaterialCategory, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE name ILIKE $1 AND company_id = $2`,
		domain.TableMaterialCategories)

	searchQuery := param.Query + "%"

	err := mc.psql.QueryRowContext(ctx, countQuery, searchQuery, param.CompanyId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT 
		    id, name, company_id, description, slug, created_at, updated_at, is_active, img_url 
		FROM %s WHERE name ILIKE $1 AND company_id = $2 ORDER BY %s %s LIMIT $3 OFFSET $4`,
		domain.TableMaterialCategories, param.SortField, param.Sort)

	rows, err := mc.psql.QueryContext(ctx, query, searchQuery, param.CompanyId, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var categories []domain.MaterialCategory

	for rows.Next() {
		var c domain.MaterialCategory
		if err = rows.Scan(
			&c.ID, &c.Name, &c.CompanyID, &c.Description, &c.Slug, &c.CreatedAt, &c.UpdatedAt, &c.IsActive, &c.ImgURL,
		); err != nil {
			return nil, 0, err
		}

		categories = append(categories, c)
	}

	return categories, totalCount, nil
}
