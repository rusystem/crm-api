package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/rusystem/crm-api/pkg/domain"
)

type Warehouse interface {
	Create(ctx context.Context, warehouse domain.Warehouse) (int64, error)
	GetById(ctx context.Context, id int64) (domain.Warehouse, error)
	Update(ctx context.Context, warehouse domain.Warehouse) error
	Delete(ctx context.Context, id int64) error
	GetListByCompanyId(ctx context.Context, id int64, param domain.Param) ([]domain.Warehouse, int64, error)
	GetResponsibleUsers(ctx context.Context, companyId int64, param domain.Param) ([]domain.User, int64, error)
}

type WarehousePostgresRepository struct {
	db *sql.DB
}

func NewWarehousePostgresRepository(db *sql.DB) *WarehousePostgresRepository {
	return &WarehousePostgresRepository{db: db}
}

func (wpr *WarehousePostgresRepository) Create(ctx context.Context, warehouse domain.Warehouse) (int64, error) {
	otherFieldsJSON, err := json.Marshal(warehouse.OtherFields)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	query := fmt.Sprintf(`
    INSERT INTO %s (
        name, address, responsible_person, phone, email,
        max_capacity, current_occupancy, other_fields, country, region,
        comments, company_id
    ) VALUES (
        $1, $2, $3, $4, $5,
        $6, $7, $8, $9, $10,
        $11, $12
    ) RETURNING id
    `, domain.TableWarehouse)

	var id int64
	if err = wpr.db.QueryRowContext(ctx, query,
		warehouse.Name, warehouse.Address, warehouse.ResponsiblePerson, warehouse.Phone, warehouse.Email,
		warehouse.MaxCapacity, warehouse.CurrentOccupancy, otherFieldsJSON, warehouse.Country, warehouse.Region,
		warehouse.Comments, warehouse.CompanyId,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert warehouse: %v", err)
	}

	return id, nil
}

func (wpr *WarehousePostgresRepository) GetById(ctx context.Context, id int64) (domain.Warehouse, error) {
	query := fmt.Sprintf(`
    SELECT
        id, name, address, responsible_person, phone, email,
        max_capacity, current_occupancy, other_fields, country, region,
        comments, created_at, company_id
    FROM %s
    WHERE id = $1;
    `, domain.TableWarehouse)

	var warehouse domain.Warehouse
	var otherFieldsJSON []byte

	row := wpr.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&warehouse.ID, &warehouse.Name, &warehouse.Address, &warehouse.ResponsiblePerson, &warehouse.Phone, &warehouse.Email,
		&warehouse.MaxCapacity, &warehouse.CurrentOccupancy, &otherFieldsJSON, &warehouse.Country, &warehouse.Region,
		&warehouse.Comments, &warehouse.CreatedAt, &warehouse.CompanyId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Warehouse{}, domain.ErrWarehouseNotFound
		}

		return domain.Warehouse{}, err
	}

	if err = json.Unmarshal(otherFieldsJSON, &warehouse.OtherFields); err != nil {
		return domain.Warehouse{}, fmt.Errorf("failed to unmarshal other_fields JSON: %v", err)
	}

	return warehouse, nil
}

func (wpr *WarehousePostgresRepository) Update(ctx context.Context, warehouse domain.Warehouse) error {
	otherFieldsJSON, err := json.Marshal(warehouse.OtherFields)
	if err != nil {
		return fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	query := fmt.Sprintf(`
	UPDATE %s
	SET
		name = $1, address = $2, responsible_person = $3, phone = $4, email = $5,
		max_capacity = $6, current_occupancy = $7, other_fields = $8, country = $9,
		region = $10, comments = $11
	WHERE id = $12
	`, domain.TableWarehouse)

	_, err = wpr.db.ExecContext(ctx, query,
		warehouse.Name, warehouse.Address, warehouse.ResponsiblePerson, warehouse.Phone, warehouse.Email,
		warehouse.MaxCapacity, warehouse.CurrentOccupancy, otherFieldsJSON, warehouse.Country,
		warehouse.Region, warehouse.Comments,
		warehouse.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update warehouse: %v", err)
	}

	return nil
}

func (wpr *WarehousePostgresRepository) Delete(ctx context.Context, id int64) error {
	_, err := wpr.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", domain.TableWarehouse), id)
	return err
}

func (wpr *WarehousePostgresRepository) GetListByCompanyId(ctx context.Context, id int64, param domain.Param) ([]domain.Warehouse, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
	SELECT COUNT(*)
	FROM %s
	WHERE company_id = $1
	`, domain.TableWarehouse)

	err := wpr.db.QueryRowContext(ctx, countQuery, id).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count warehouses by company ID: %v", err)
	}

	query := fmt.Sprintf(`
	SELECT
		id, name, address, responsible_person, phone, email,
		max_capacity, current_occupancy, other_fields, country, region, 
		comments, created_at, company_id
	FROM %s
	WHERE company_id = $1 ORDER BY %s %s
	LIMIT $2 OFFSET $3;
	`, domain.TableWarehouse, param.SortField, param.Sort)

	rows, err := wpr.db.QueryContext(ctx, query, id, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get warehouses by company ID: %v", err)
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var warehouses []domain.Warehouse
	for rows.Next() {
		var warehouse domain.Warehouse
		var otherFieldsJSON []byte
		if err = rows.Scan(
			&warehouse.ID, &warehouse.Name, &warehouse.Address, &warehouse.ResponsiblePerson, &warehouse.Phone, &warehouse.Email,
			&warehouse.MaxCapacity, &warehouse.CurrentOccupancy, &otherFieldsJSON, &warehouse.Country, &warehouse.Region,
			&warehouse.Comments, &warehouse.CreatedAt, &warehouse.CompanyId,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan warehouse: %v", err)
		}

		if err = json.Unmarshal(otherFieldsJSON, &warehouse.OtherFields); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal other_fields JSON: %v", err)
		}

		warehouses = append(warehouses, warehouse)
	}

	return warehouses, totalCount, nil
}

func (wpr *WarehousePostgresRepository) GetResponsibleUsers(ctx context.Context, companyId int64, param domain.Param) ([]domain.User, int64, error) {
	var users []domain.User
	var totalCount int64

	sections := pq.Array([]string{domain.SectionFullCompanyAccess, domain.SectionFullAccess, domain.SectionPurchasePlanningAccess})

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE company_id = $1 AND EXISTS (
		    SELECT 1 FROM jsonb_array_elements_text(sections) AS section
		    WHERE section = ANY ($2)
		)
	`, domain.UsersTable)

	if err := wpr.db.QueryRowContext(ctx, countQuery, companyId, sections).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("failed to count responsible users: %v", err)
	}

	query := fmt.Sprintf(`
		SELECT 
		    id, company_id, username, name, email, phone, password_hash, created_at, 
		    updated_at, last_login, is_active, role, language, country, 
		    is_approved, is_send_system_notification, sections, position
		FROM %s
		WHERE company_id = $1 AND EXISTS (SELECT 1 FROM jsonb_array_elements_text(sections) AS section
		WHERE section = ANY ($2)) ORDER BY %s %s
		LIMIT $3 OFFSET $4;
		`, domain.UsersTable, param.SortField, param.Sort)

	rows, err := wpr.db.QueryContext(ctx, query, companyId, sections, param.Limit, param.Offset)
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

		if err := json.Unmarshal(b, &user.Sections); err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}
