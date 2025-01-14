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

type Suppliers interface {
	Create(ctx context.Context, supplier domain.Supplier) (int64, error)
	GetById(ctx context.Context, id int64) (domain.Supplier, error)
	Update(ctx context.Context, supplier domain.Supplier) error
	Delete(ctx context.Context, id int64) error
	GetListByCompanyId(ctx context.Context, id int64, param domain.Param) ([]domain.Supplier, int64, error)
}

type SuppliersPostgresRepository struct {
	psql *sql.DB
}

func NewSuppliersPostgresRepository(psql *sql.DB) *SuppliersPostgresRepository {
	return &SuppliersPostgresRepository{psql: psql}
}

func (sr *SuppliersPostgresRepository) Create(ctx context.Context, supplier domain.Supplier) (int64, error) {
	otherFieldsJSON, err := json.Marshal(supplier.OtherFields)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (name, legal_address, actual_address, warehouse_address, contact_person, phone, email, 
		                       website, contract_number, product_categories, purchase_amount, balance, product_types, 
		                       comments, files, country, region, tax_id, bank_details, registration_date, payment_terms, 
		                       is_active, other_fields, company_id, contract_date, locality) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26) RETURNING id`,
		domain.TableSupplier)

	var id int64
	if err = sr.psql.QueryRowContext(ctx, query, supplier.Name, supplier.LegalAddress, supplier.ActualAddress,
		supplier.WarehouseAddress, supplier.ContactPerson, supplier.Phone, supplier.Email, supplier.Website,
		supplier.ContractNumber, pq.Array(supplier.ProductCategories), supplier.PurchaseAmount, supplier.Balance, supplier.ProductTypes,
		supplier.Comments, supplier.Files, supplier.Country, supplier.Region, supplier.TaxID, supplier.BankDetails,
		supplier.RegistrationDate, pq.Array(supplier.PaymentTerms), supplier.IsActive, otherFieldsJSON, supplier.CompanyId, supplier.ContractDate,
		supplier.Locality,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (sr *SuppliersPostgresRepository) GetById(ctx context.Context, id int64) (domain.Supplier, error) {
	query := fmt.Sprintf(`
    SELECT
        id, name, legal_address, actual_address, warehouse_address,
        contact_person, phone, email, website, contract_number,
        product_categories, purchase_amount, balance, product_types,
        comments, files, country, region, tax_id, bank_details,
        registration_date, payment_terms, is_active, other_fields, 
        company_id, contract_date, locality
    FROM %s
    WHERE id = $1;
    `, domain.TableSupplier)

	var supplier domain.Supplier
	var otherFieldsJSON []byte

	// Выполнение запроса и сканирование результата в объект Supplier
	row := sr.psql.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&supplier.ID, &supplier.Name, &supplier.LegalAddress, &supplier.ActualAddress,
		&supplier.WarehouseAddress, &supplier.ContactPerson, &supplier.Phone, &supplier.Email,
		&supplier.Website, &supplier.ContractNumber, pq.Array(&supplier.ProductCategories), &supplier.PurchaseAmount,
		&supplier.Balance, &supplier.ProductTypes, &supplier.Comments, &supplier.Files,
		&supplier.Country, &supplier.Region, &supplier.TaxID, &supplier.BankDetails,
		&supplier.RegistrationDate, pq.Array(&supplier.PaymentTerms), &supplier.IsActive, &otherFieldsJSON,
		&supplier.CompanyId, &supplier.ContractDate, &supplier.Locality,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Supplier{}, domain.ErrSupplierNotFound
		}

		return domain.Supplier{}, err
	}

	if err = json.Unmarshal(otherFieldsJSON, &supplier.OtherFields); err != nil {
		return domain.Supplier{}, fmt.Errorf("failed to unmarshal other_fields JSON: %v", err)
	}

	return supplier, nil
}

func (sr *SuppliersPostgresRepository) Update(ctx context.Context, supplier domain.Supplier) error {
	otherFieldsJSON, err := json.Marshal(supplier.OtherFields)
	if err != nil {
		return fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET name = $1, legal_address = $2, actual_address = $3, warehouse_address = $4, contact_person = $5,
			phone = $6, email = $7, website = $8, contract_number = $9, product_categories = $10, purchase_amount = $11,
			balance = $12, product_types = $13, comments = $14, files = $15, country = $16, region = $17, tax_id = $18,
			bank_details = $19, registration_date = $20, payment_terms = $21, is_active = $22, other_fields = $23, contract_date = $24,
		    locality = $25
		WHERE id = $26;
	`, domain.TableSupplier)

	_, err = sr.psql.ExecContext(ctx, query, supplier.Name, supplier.LegalAddress, supplier.ActualAddress,
		supplier.WarehouseAddress, supplier.ContactPerson, supplier.Phone, supplier.Email, supplier.Website,
		supplier.ContractNumber, pq.Array(supplier.ProductCategories), supplier.PurchaseAmount, supplier.Balance,
		supplier.ProductTypes, supplier.Comments, supplier.Files, supplier.Country, supplier.Region, supplier.TaxID,
		supplier.BankDetails, supplier.RegistrationDate, pq.Array(supplier.PaymentTerms), supplier.IsActive, otherFieldsJSON,
		supplier.ContractDate, supplier.Locality, supplier.ID)

	return err
}

func (sr *SuppliersPostgresRepository) Delete(ctx context.Context, id int64) error {
	_, err := sr.psql.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", domain.TableSupplier), id)
	return err
}

func (sr *SuppliersPostgresRepository) GetListByCompanyId(ctx context.Context, id int64, param domain.Param) ([]domain.Supplier, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
	SELECT COUNT(*)
	FROM %s
	WHERE company_id = $1
	`, domain.TableSupplier)

	err := sr.psql.QueryRowContext(ctx, countQuery, id).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
	SELECT
		id, name, legal_address, actual_address, warehouse_address,
		contact_person, phone, email, website, contract_number,
		product_categories, purchase_amount, balance, product_types,
		comments, files, country, region, tax_id, bank_details,
		registration_date, payment_terms, is_active, other_fields, company_id, 
		contract_date, locality
	FROM %s
	WHERE company_id = $1 ORDER BY %s %s
	LIMIT $2 OFFSET $3;
	`, domain.TableSupplier, param.SortField, param.Sort)

	rows, err := sr.psql.QueryContext(ctx, query, id, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var suppliers []domain.Supplier

	for rows.Next() {
		var supplier domain.Supplier
		var otherFieldsJSON []byte

		if err = rows.Scan(
			&supplier.ID, &supplier.Name, &supplier.LegalAddress, &supplier.ActualAddress,
			&supplier.WarehouseAddress, &supplier.ContactPerson, &supplier.Phone, &supplier.Email,
			&supplier.Website, &supplier.ContractNumber, pq.Array(&supplier.ProductCategories), &supplier.PurchaseAmount,
			&supplier.Balance, &supplier.ProductTypes, &supplier.Comments, &supplier.Files,
			&supplier.Country, &supplier.Region, &supplier.TaxID, &supplier.BankDetails,
			&supplier.RegistrationDate, pq.Array(&supplier.PaymentTerms), &supplier.IsActive, &otherFieldsJSON,
			&supplier.CompanyId, &supplier.ContractDate, &supplier.Locality,
		); err != nil {
			return nil, 0, err
		}

		if err = json.Unmarshal(otherFieldsJSON, &supplier.OtherFields); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal other_fields JSON: %v", err)
		}

		suppliers = append(suppliers, supplier)
	}

	return suppliers, totalCount, nil
}
