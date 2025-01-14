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

type Materials interface {
	CreatePlanning(ctx context.Context, material domain.Material) (int64, error)
	UpdatePlanning(ctx context.Context, material domain.Material) error
	DeletePlanning(ctx context.Context, id int64) error
	GetPlanningById(ctx context.Context, id int64) (domain.Material, error)
	GetPlanningList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error)
	MovePlanningToPurchased(ctx context.Context, id int64) (int64, int64, error)

	CreatePurchased(ctx context.Context, material domain.Material) (int64, int64, error)
	UpdatePurchased(ctx context.Context, material domain.Material) error
	DeletePurchased(ctx context.Context, id int64) error
	GetPurchasedById(ctx context.Context, id int64) (domain.Material, error)
	GetPurchasedList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error)
	MovePurchasedToArchive(ctx context.Context, id int64) error

	GetPlanningArchiveById(ctx context.Context, id int64) (domain.Material, error)
	GetPurchasedArchiveById(ctx context.Context, id int64) (domain.Material, error)
	GetPlanningArchiveList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error)
	GetPurchasedArchiveList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error)
	DeletePlanningArchive(ctx context.Context, id int64) error
	DeletePurchasedArchive(ctx context.Context, id int64) error

	Search(ctx context.Context, param domain.MaterialParams) ([]domain.Material, int64, error)

	GetIncomeHistoryByWarehouseId(ctx context.Context, id int64, param domain.Param) ([]domain.Material, int64, error)
}

type MaterialsPostgresRepository struct {
	psql *sql.DB
}

func NewMaterialsPostgresRepository(psql *sql.DB) *MaterialsPostgresRepository {
	return &MaterialsPostgresRepository{
		psql: psql,
	}
}

func (mr *MaterialsPostgresRepository) CreatePlanning(ctx context.Context, material domain.Material) (int64, error) {
	otherFieldsJSON, err := json.Marshal(material.OtherFields)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume, 
						price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve, 
						received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost, 
						warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, supplier_name, contract_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, 
				$22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32) RETURNING id`,
		domain.TablePlanningMaterials)

	var id int64
	if err = mr.psql.QueryRowContext(ctx, query,
		material.WarehouseID, material.ItemID, material.Name, material.ByInvoice, material.Article, pq.Array(material.ProductCategory),
		material.Unit, material.TotalQuantity, material.Volume, material.PriceWithoutVAT, material.TotalWithoutVAT,
		material.SupplierID, material.Location, material.ContractDate, material.File, material.Status, material.Comments,
		material.Reserve, material.ReceivedDate, material.LastUpdated, material.MinStockLevel, material.ExpirationDate,
		material.ResponsiblePerson, material.StorageCost, material.WarehouseSection, material.IncomingDeliveryNumber,
		otherFieldsJSON, material.CompanyID, material.InternalName, material.UnitsPerPackage, material.SupplierName,
		material.ContractNumber,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert planning material: %v", err)
	}

	return id, nil
}

func (mr *MaterialsPostgresRepository) UpdatePlanning(ctx context.Context, material domain.Material) error {
	otherFieldsJSON, err := json.Marshal(material.OtherFields)
	if err != nil {
		return fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET
			warehouse_id = $1, item_id = $2, name = $3, by_invoice = $4, article = $5, product_category = $6, unit = $7,
			total_quantity = $8, volume = $9, price_without_vat = $10, total_without_vat = $11, supplier_id = $12, location = $13,
			contract_date = $14, file = $15, status = $16, comments = $17, reserve = $18, received_date = $19, last_updated = $20,
			min_stock_level = $21, expiration_date = $22, responsible_person = $23, storage_cost = $24, warehouse_section = $25,
			incoming_delivery_number = $26, other_fields = $27, internal_name = $28, units_per_package = $29, supplier_name = $30, 
		    contract_number = $31
		WHERE id = $32`,
		domain.TablePlanningMaterials)

	_, err = mr.psql.ExecContext(ctx, query,
		material.WarehouseID, material.ItemID, material.Name, material.ByInvoice, material.Article, pq.Array(material.ProductCategory),
		material.Unit, material.TotalQuantity, material.Volume, material.PriceWithoutVAT, material.TotalWithoutVAT,
		material.SupplierID, material.Location, material.ContractDate, material.File, material.Status, material.Comments,
		material.Reserve, material.ReceivedDate, material.LastUpdated, material.MinStockLevel, material.ExpirationDate,
		material.ResponsiblePerson, material.StorageCost, material.WarehouseSection, material.IncomingDeliveryNumber, otherFieldsJSON,
		material.InternalName, material.UnitsPerPackage, material.SupplierName, material.ContractNumber, material.ID,
	)

	return err
}

func (mr *MaterialsPostgresRepository) DeletePlanning(ctx context.Context, id int64) error {
	_, err := mr.psql.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", domain.TablePlanningMaterials), id)
	return err
}

func (mr *MaterialsPostgresRepository) GetPlanningById(ctx context.Context, id int64) (domain.Material, error) {
	return mr.getPlanningById(ctx, id)
}

func (mr *MaterialsPostgresRepository) getPlanningById(ctx context.Context, id int64) (domain.Material, error) {
	query := fmt.Sprintf(`
	SELECT 
	    id, warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume,
		price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve,
		received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost,
		warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, supplier_name, 
		contract_number
	FROM %s WHERE id = $1
	`, domain.TablePlanningMaterials)

	var material domain.Material
	var otherFieldsJSON []byte

	if err := mr.psql.QueryRowContext(ctx, query, id).Scan(
		&material.ID, &material.WarehouseID, &material.ItemID, &material.Name, &material.ByInvoice, &material.Article,
		pq.Array(&material.ProductCategory), &material.Unit, &material.TotalQuantity, &material.Volume,
		&material.PriceWithoutVAT, &material.TotalWithoutVAT, &material.SupplierID, &material.Location,
		&material.ContractDate, &material.File, &material.Status, &material.Comments, &material.Reserve,
		&material.ReceivedDate, &material.LastUpdated, &material.MinStockLevel, &material.ExpirationDate,
		&material.ResponsiblePerson, &material.StorageCost, &material.WarehouseSection,
		&material.IncomingDeliveryNumber, &otherFieldsJSON, &material.CompanyID, &material.InternalName,
		&material.UnitsPerPackage, &material.SupplierName, &material.ContractNumber,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Material{}, domain.ErrMaterialNotFound
		}

		return domain.Material{}, err
	}

	if err := json.Unmarshal(otherFieldsJSON, &material.OtherFields); err != nil {
		return domain.Material{}, err
	}

	return material, nil
}

func (mr *MaterialsPostgresRepository) GetPlanningList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
	SELECT COUNT(*)
	FROM %s
	WHERE company_id = $1
	`, domain.TablePlanningMaterials)

	err := mr.psql.QueryRowContext(ctx, countQuery, params.CompanyId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
	SELECT 
	    id, warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume,
		price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve,
		received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost,
		warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		supplier_name, contract_number
	FROM %s WHERE company_id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3
	`, domain.TablePlanningMaterials, params.SortField, params.Sort)

	rows, err := mr.psql.QueryContext(ctx, query, params.CompanyId, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var materials []domain.Material

	for rows.Next() {
		var material domain.Material
		var otherFieldsJSON []byte

		if err = rows.Scan(
			&material.ID, &material.WarehouseID, &material.ItemID, &material.Name, &material.ByInvoice, &material.Article,
			pq.Array(&material.ProductCategory), &material.Unit, &material.TotalQuantity, &material.Volume,
			&material.PriceWithoutVAT, &material.TotalWithoutVAT, &material.SupplierID, &material.Location,
			&material.ContractDate, &material.File, &material.Status, &material.Comments, &material.Reserve,
			&material.ReceivedDate, &material.LastUpdated, &material.MinStockLevel, &material.ExpirationDate,
			&material.ResponsiblePerson, &material.StorageCost, &material.WarehouseSection,
			&material.IncomingDeliveryNumber, &otherFieldsJSON, &material.CompanyID, &material.InternalName,
			&material.UnitsPerPackage, &material.SupplierName, &material.ContractNumber,
		); err != nil {
			return nil, 0, err
		}

		if err = json.Unmarshal(otherFieldsJSON, &material.OtherFields); err != nil {
			return nil, 0, err
		}

		materials = append(materials, material)
	}

	return materials, totalCount, nil
}

func (mr *MaterialsPostgresRepository) MovePlanningToPurchased(ctx context.Context, id int64) (int64, int64, error) {
	// Удаляем из planning, переносим сразу в purchased и archived
	tx, err := mr.psql.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer func(tx *sql.Tx) {
		if err = tx.Rollback(); err != nil {
			return
		}
	}(tx)

	material, err := mr.getPlanningById(ctx, id)
	if err != nil {
		return 0, 0, err
	}

	otherFieldsJSON, err := json.Marshal(material.OtherFields)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	// 1. удаляем из planning
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1",
		domain.TablePlanningMaterials)

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return 0, 0, err
	}

	// 2. переносим в purchased
	query = fmt.Sprintf(`
		INSERT INTO %s (warehouse_id, name, by_invoice, article, product_category, unit, total_quantity, volume, 
						price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve, 
						received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost, 
						warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		                supplier_name, contract_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, 
				$22, $23, $24, $25, $26, $27, $28, $29, $30, $31) RETURNING id, item_id`,
		domain.TablePurchasedMaterials)

	var newId int64
	var itemId int64
	if err = tx.QueryRowContext(ctx, query,
		material.WarehouseID, material.Name, material.ByInvoice, material.Article, pq.Array(material.ProductCategory),
		material.Unit, material.TotalQuantity, material.Volume, material.PriceWithoutVAT, material.TotalWithoutVAT,
		material.SupplierID, material.Location, material.ContractDate, material.File, material.Status, material.Comments,
		material.Reserve, material.ReceivedDate, material.LastUpdated, material.MinStockLevel, material.ExpirationDate,
		material.ResponsiblePerson, material.StorageCost, material.WarehouseSection, material.IncomingDeliveryNumber,
		otherFieldsJSON, material.CompanyID, material.InternalName, material.UnitsPerPackage, material.SupplierName,
		material.ContractNumber,
	).Scan(&newId, &itemId); err != nil {
		return 0, 0, fmt.Errorf("failed to insert purchased material: %v", err)
	}

	// 3. переносим в planning archive
	query = fmt.Sprintf(`
		INSERT INTO %s (warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume, 
						price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve, 
						received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost, 
						warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		                supplier_name, contract_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, 
				$22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32)`,
		domain.TablePlanningMaterialsArchive)

	_, err = tx.ExecContext(ctx, query,
		material.WarehouseID, material.ItemID, material.Name, material.ByInvoice, material.Article, pq.Array(material.ProductCategory),
		material.Unit, material.TotalQuantity, material.Volume, material.PriceWithoutVAT, material.TotalWithoutVAT,
		material.SupplierID, material.Location, material.ContractDate, material.File, material.Status, material.Comments,
		material.Reserve, material.ReceivedDate, material.LastUpdated, material.MinStockLevel, material.ExpirationDate,
		material.ResponsiblePerson, material.StorageCost, material.WarehouseSection, material.IncomingDeliveryNumber,
		otherFieldsJSON, material.CompanyID, material.InternalName, material.UnitsPerPackage, material.SupplierName,
		material.ContractNumber,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to insert purchased archive material: %v", err)
	}

	return newId, itemId, tx.Commit()
}

func (mr *MaterialsPostgresRepository) CreatePurchased(ctx context.Context, material domain.Material) (int64, int64, error) {
	otherFieldsJSON, err := json.Marshal(material.OtherFields)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (warehouse_id, name, by_invoice, article, product_category, unit, total_quantity, volume, 
						price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve, 
						received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost, 
						warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		                supplier_name, contract_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, 
				$22, $23, $24, $25, $26, $27, $28, $29, $30, $31) RETURNING id, item_id`,
		domain.TablePurchasedMaterials)

	var id int64
	var itemId int64
	if err = mr.psql.QueryRowContext(ctx, query,
		material.WarehouseID, material.Name, material.ByInvoice, material.Article, pq.Array(material.ProductCategory),
		material.Unit, material.TotalQuantity, material.Volume, material.PriceWithoutVAT, material.TotalWithoutVAT,
		material.SupplierID, material.Location, material.ContractDate, material.File, material.Status, material.Comments,
		material.Reserve, material.ReceivedDate, material.LastUpdated, material.MinStockLevel, material.ExpirationDate,
		material.ResponsiblePerson, material.StorageCost, material.WarehouseSection, material.IncomingDeliveryNumber,
		otherFieldsJSON, material.CompanyID, material.InternalName, material.UnitsPerPackage, material.SupplierName,
		material.ContractNumber,
	).Scan(&id, &itemId); err != nil {
		return 0, 0, fmt.Errorf("failed to insert purchased material: %v", err)
	}

	return id, itemId, nil
}

func (mr *MaterialsPostgresRepository) UpdatePurchased(ctx context.Context, material domain.Material) error {
	otherFieldsJSON, err := json.Marshal(material.OtherFields)
	if err != nil {
		return fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET
			warehouse_id = $1, item_id = $2, name = $3, by_invoice = $4, article = $5, product_category = $6, unit = $7,
			total_quantity = $8, volume = $9, price_without_vat = $10, total_without_vat = $11, supplier_id = $12, location = $13,
			contract_date = $14, file = $15, status = $16, comments = $17, reserve = $18, received_date = $19, last_updated = $20,
			min_stock_level = $21, expiration_date = $22, responsible_person = $23, storage_cost = $24, warehouse_section = $25,
			incoming_delivery_number = $26, other_fields = $27, internal_name = $28, units_per_package = $29, supplier_name = $30, 
		    contract_number = $31
		WHERE id = $32`,
		domain.TablePurchasedMaterials)

	_, err = mr.psql.ExecContext(ctx, query,
		material.WarehouseID, material.ItemID, material.Name, material.ByInvoice, material.Article, pq.Array(material.ProductCategory),
		material.Unit, material.TotalQuantity, material.Volume, material.PriceWithoutVAT, material.TotalWithoutVAT,
		material.SupplierID, material.Location, material.ContractDate, material.File, material.Status, material.Comments,
		material.Reserve, material.ReceivedDate, material.LastUpdated, material.MinStockLevel, material.ExpirationDate,
		material.ResponsiblePerson, material.StorageCost, material.WarehouseSection, material.IncomingDeliveryNumber,
		otherFieldsJSON, material.InternalName, material.UnitsPerPackage, material.SupplierName,
		material.ContractNumber, material.ID,
	)

	return err
}

func (mr *MaterialsPostgresRepository) DeletePurchased(ctx context.Context, id int64) error {
	_, err := mr.psql.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", domain.TablePurchasedMaterials), id)
	return err
}

func (mr *MaterialsPostgresRepository) GetPurchasedById(ctx context.Context, id int64) (domain.Material, error) {
	return mr.getPurchasedById(ctx, id)
}

func (mr *MaterialsPostgresRepository) getPurchasedById(ctx context.Context, id int64) (domain.Material, error) {
	query := fmt.Sprintf(`
	SELECT 
	    id, warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume,
		price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve,
		received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost,
		warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		supplier_name, contract_number
	FROM %s WHERE id = $1
	`, domain.TablePurchasedMaterials)

	var material domain.Material
	var otherFieldsJSON []byte

	if err := mr.psql.QueryRowContext(ctx, query, id).Scan(
		&material.ID, &material.WarehouseID, &material.ItemID, &material.Name, &material.ByInvoice, &material.Article,
		pq.Array(&material.ProductCategory), &material.Unit, &material.TotalQuantity, &material.Volume,
		&material.PriceWithoutVAT, &material.TotalWithoutVAT, &material.SupplierID, &material.Location,
		&material.ContractDate, &material.File, &material.Status, &material.Comments, &material.Reserve,
		&material.ReceivedDate, &material.LastUpdated, &material.MinStockLevel, &material.ExpirationDate,
		&material.ResponsiblePerson, &material.StorageCost, &material.WarehouseSection, &material.IncomingDeliveryNumber,
		&otherFieldsJSON, &material.CompanyID, &material.InternalName, &material.UnitsPerPackage, &material.SupplierName,
		&material.ContractNumber,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Material{}, domain.ErrMaterialNotFound
		}

		return domain.Material{}, err
	}

	if err := json.Unmarshal(otherFieldsJSON, &material.OtherFields); err != nil {
		return domain.Material{}, err
	}

	return material, nil
}

func (mr *MaterialsPostgresRepository) GetPurchasedList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error) {
	var totalCount int64

	// Подсчет общего количества записей
	countQuery := fmt.Sprintf(`
	SELECT COUNT(*)
	FROM %s
	WHERE company_id = $1
	`, domain.TablePurchasedMaterials)

	err := mr.psql.QueryRowContext(ctx, countQuery, params.CompanyId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
	SELECT 
	    id, warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume,
		price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve,
		received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost,
		warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		supplier_name, contract_number
	FROM %s WHERE company_id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3
	`, domain.TablePurchasedMaterials, params.SortField, params.Sort)

	rows, err := mr.psql.QueryContext(ctx, query, params.CompanyId, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var materials []domain.Material

	for rows.Next() {
		var material domain.Material
		var otherFieldsJSON []byte

		if err = rows.Scan(
			&material.ID, &material.WarehouseID, &material.ItemID, &material.Name, &material.ByInvoice, &material.Article,
			pq.Array(&material.ProductCategory), &material.Unit, &material.TotalQuantity, &material.Volume,
			&material.PriceWithoutVAT, &material.TotalWithoutVAT, &material.SupplierID, &material.Location,
			&material.ContractDate, &material.File, &material.Status, &material.Comments, &material.Reserve,
			&material.ReceivedDate, &material.LastUpdated, &material.MinStockLevel, &material.ExpirationDate,
			&material.ResponsiblePerson, &material.StorageCost, &material.WarehouseSection,
			&material.IncomingDeliveryNumber, &otherFieldsJSON, &material.CompanyID, &material.InternalName,
			&material.UnitsPerPackage, &material.SupplierName, &material.ContractNumber,
		); err != nil {
			return nil, 0, err
		}

		if err = json.Unmarshal(otherFieldsJSON, &material.OtherFields); err != nil {
			return nil, 0, err
		}

		materials = append(materials, material)
	}

	return materials, totalCount, nil
}

func (mr *MaterialsPostgresRepository) MovePurchasedToArchive(ctx context.Context, id int64) error {
	// Удаляем из purchased и переносим сразу в archived
	tx, err := mr.psql.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		if err = tx.Rollback(); err != nil {
			return
		}
	}(tx)

	material, err := mr.getPurchasedById(ctx, id)
	if err != nil {
		return err
	}

	otherFieldsJSON, err := json.Marshal(material.OtherFields)
	if err != nil {
		return fmt.Errorf("failed to marshal other_fields to JSON: %v", err)
	}

	// 1. удаляем из purchased
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1",
		domain.TablePurchasedMaterials)

	if _, err = tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	// 2. переносим в purchased archive
	query = fmt.Sprintf(`
		INSERT INTO %s (warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume, 
						price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve, 
						received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost, 
						warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		                supplier_name, contract_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, 
				$22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32)`,
		domain.TablePurchasedMaterialsArchive)

	_, err = tx.ExecContext(ctx, query,
		material.WarehouseID, material.ItemID, material.Name, material.ByInvoice, material.Article, pq.Array(material.ProductCategory),
		material.Unit, material.TotalQuantity, material.Volume, material.PriceWithoutVAT, material.TotalWithoutVAT,
		material.SupplierID, material.Location, material.ContractDate, material.File, material.Status, material.Comments,
		material.Reserve, material.ReceivedDate, material.LastUpdated, material.MinStockLevel, material.ExpirationDate,
		material.ResponsiblePerson, material.StorageCost, material.WarehouseSection, material.IncomingDeliveryNumber,
		otherFieldsJSON, material.CompanyID, material.InternalName, material.UnitsPerPackage, material.SupplierName,
		material.ContractNumber,
	)
	if err != nil {
		return fmt.Errorf("failed to insert purchased material archive: %v", err)
	}

	return tx.Commit()
}

func (mr *MaterialsPostgresRepository) GetPlanningArchiveById(ctx context.Context, id int64) (domain.Material, error) {
	query := fmt.Sprintf(`
	SELECT 
	    id, warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume,
		price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve,
		received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost,
		warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		supplier_name, contract_number
	FROM %s WHERE id = $1
	`, domain.TablePlanningMaterialsArchive)

	var material domain.Material
	var otherFieldsJSON []byte

	if err := mr.psql.QueryRowContext(ctx, query, id).Scan(
		&material.ID, &material.WarehouseID, &material.ItemID, &material.Name, &material.ByInvoice, &material.Article,
		pq.Array(&material.ProductCategory), &material.Unit, &material.TotalQuantity, &material.Volume,
		&material.PriceWithoutVAT, &material.TotalWithoutVAT, &material.SupplierID, &material.Location,
		&material.ContractDate, &material.File, &material.Status, &material.Comments, &material.Reserve,
		&material.ReceivedDate, &material.LastUpdated, &material.MinStockLevel, &material.ExpirationDate,
		&material.ResponsiblePerson, &material.StorageCost, &material.WarehouseSection,
		&material.IncomingDeliveryNumber, &otherFieldsJSON, &material.CompanyID, &material.InternalName,
		&material.UnitsPerPackage, &material.SupplierName, &material.ContractNumber,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Material{}, domain.ErrMaterialNotFound
		}

		return domain.Material{}, err
	}

	if err := json.Unmarshal(otherFieldsJSON, &material.OtherFields); err != nil {
		return domain.Material{}, err
	}

	return material, nil
}

func (mr *MaterialsPostgresRepository) GetPurchasedArchiveById(ctx context.Context, id int64) (domain.Material, error) {
	query := fmt.Sprintf(`
	SELECT 
	    id, warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume,
		price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve,
		received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost,
		warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		supplier_name, contract_number
	FROM %s WHERE id = $1
	`, domain.TablePurchasedMaterialsArchive)

	var material domain.Material
	var otherFieldsJSON []byte

	if err := mr.psql.QueryRowContext(ctx, query, id).Scan(
		&material.ID, &material.WarehouseID, &material.ItemID, &material.Name, &material.ByInvoice, &material.Article,
		pq.Array(&material.ProductCategory), &material.Unit, &material.TotalQuantity, &material.Volume,
		&material.PriceWithoutVAT, &material.TotalWithoutVAT, &material.SupplierID, &material.Location,
		&material.ContractDate, &material.File, &material.Status, &material.Comments, &material.Reserve,
		&material.ReceivedDate, &material.LastUpdated, &material.MinStockLevel, &material.ExpirationDate,
		&material.ResponsiblePerson, &material.StorageCost, &material.WarehouseSection,
		&material.IncomingDeliveryNumber, &otherFieldsJSON, &material.CompanyID, &material.InternalName,
		&material.UnitsPerPackage, &material.SupplierName, &material.ContractNumber,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Material{}, domain.ErrMaterialNotFound
		}

		return domain.Material{}, err
	}

	if err := json.Unmarshal(otherFieldsJSON, &material.OtherFields); err != nil {
		return domain.Material{}, err
	}

	return material, nil
}

func (mr *MaterialsPostgresRepository) GetPlanningArchiveList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
	SELECT COUNT(*)
	FROM %s
	WHERE company_id = $1
	`, domain.TablePlanningMaterialsArchive)

	err := mr.psql.QueryRowContext(ctx, countQuery, params.CompanyId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
	SELECT 
	    id, warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume,
		price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve,
		received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost,
		warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		supplier_name, contract_number
	FROM %s WHERE company_id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3
	`, domain.TablePlanningMaterialsArchive, params.SortField, params.Sort)

	rows, err := mr.psql.QueryContext(ctx, query, params.CompanyId, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var materials []domain.Material

	for rows.Next() {
		var material domain.Material
		var otherFieldsJSON []byte

		if err = rows.Scan(
			&material.ID, &material.WarehouseID, &material.ItemID, &material.Name, &material.ByInvoice, &material.Article,
			pq.Array(&material.ProductCategory), &material.Unit, &material.TotalQuantity, &material.Volume,
			&material.PriceWithoutVAT, &material.TotalWithoutVAT, &material.SupplierID, &material.Location,
			&material.ContractDate, &material.File, &material.Status, &material.Comments, &material.Reserve,
			&material.ReceivedDate, &material.LastUpdated, &material.MinStockLevel, &material.ExpirationDate,
			&material.ResponsiblePerson, &material.StorageCost, &material.WarehouseSection,
			&material.IncomingDeliveryNumber, &otherFieldsJSON, &material.CompanyID, &material.InternalName,
			&material.UnitsPerPackage, &material.SupplierName, &material.ContractNumber,
		); err != nil {
			return nil, 0, err
		}

		if err = json.Unmarshal(otherFieldsJSON, &material.OtherFields); err != nil {
			return nil, 0, err
		}

		materials = append(materials, material)
	}

	return materials, totalCount, nil
}

func (mr *MaterialsPostgresRepository) GetPurchasedArchiveList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
	SELECT COUNT(*)
	FROM %s
	WHERE company_id = $1
	`, domain.TablePurchasedMaterialsArchive)

	err := mr.psql.QueryRowContext(ctx, countQuery, params.CompanyId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
	SELECT 
	    id, warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume,
		price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve,
		received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost,
		warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		supplier_name, contract_number
	FROM %s WHERE company_id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3
	`, domain.TablePurchasedMaterialsArchive, params.SortField, params.Sort)

	rows, err := mr.psql.QueryContext(ctx, query, params.CompanyId, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var materials []domain.Material

	for rows.Next() {
		var material domain.Material
		var otherFieldsJSON []byte

		if err = rows.Scan(
			&material.ID, &material.WarehouseID, &material.ItemID, &material.Name, &material.ByInvoice, &material.Article,
			pq.Array(&material.ProductCategory), &material.Unit, &material.TotalQuantity, &material.Volume,
			&material.PriceWithoutVAT, &material.TotalWithoutVAT, &material.SupplierID, &material.Location,
			&material.ContractDate, &material.File, &material.Status, &material.Comments, &material.Reserve,
			&material.ReceivedDate, &material.LastUpdated, &material.MinStockLevel, &material.ExpirationDate,
			&material.ResponsiblePerson, &material.StorageCost, &material.WarehouseSection,
			&material.IncomingDeliveryNumber, &otherFieldsJSON, &material.CompanyID, &material.InternalName,
			&material.UnitsPerPackage, &material.SupplierName, &material.ContractNumber,
		); err != nil {
			return nil, 0, err
		}

		if err = json.Unmarshal(otherFieldsJSON, &material.OtherFields); err != nil {
			return nil, 0, err
		}

		materials = append(materials, material)
	}

	return materials, totalCount, nil
}

func (mr *MaterialsPostgresRepository) DeletePlanningArchive(ctx context.Context, id int64) error {
	_, err := mr.psql.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", domain.TablePlanningMaterialsArchive), id)
	return err
}

func (mr *MaterialsPostgresRepository) DeletePurchasedArchive(ctx context.Context, id int64) error {
	_, err := mr.psql.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", domain.TablePurchasedMaterialsArchive), id)
	return err
}

func (mr *MaterialsPostgresRepository) Search(ctx context.Context, param domain.MaterialParams) ([]domain.Material, int64, error) {
	countQuery := fmt.Sprintf(`
        SELECT COUNT(*) FROM (
            SELECT id FROM %s WHERE name ILIKE $1 AND company_id = $2
            UNION ALL
            SELECT id FROM %s WHERE name ILIKE $1 AND company_id = $2
            UNION ALL
            SELECT id FROM %s WHERE name ILIKE $1 AND company_id = $2
            UNION ALL
            SELECT id FROM %s WHERE name ILIKE $1 AND company_id = $2
        ) AS total_count;
    `, domain.TablePlanningMaterials, domain.TablePurchasedMaterials, domain.TablePlanningMaterialsArchive, domain.TablePurchasedMaterialsArchive)

	searchPattern := param.Query + "%"

	var totalCount int64

	if err := mr.psql.QueryRowContext(ctx, countQuery, searchPattern, param.CompanyId).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	sqlQuery := fmt.Sprintf(`
        SELECT 'planning_materials' AS table_name, id, warehouse_id, name, product_category, unit, total_quantity, status, company_id 
        FROM %s
        WHERE name ILIKE $1 AND company_id = $2

        UNION ALL

        SELECT 'purchased_materials' AS table_name, id, warehouse_id, name, product_category, unit, total_quantity, status, company_id
        FROM %s
        WHERE name ILIKE $1 AND company_id = $2

        UNION ALL

        SELECT 'planning_materials_archive' AS table_name, id, warehouse_id, name, product_category, unit, total_quantity, status, company_id
        FROM %s
        WHERE name ILIKE $1 AND company_id = $2

        UNION ALL

        SELECT 'purchased_materials_archive' AS table_name, id, warehouse_id, name, product_category, unit, total_quantity, status, company_id
        FROM %s
        WHERE name ILIKE $1 AND company_id = $2

        ORDER BY %s %s
        LIMIT $3 OFFSET $4;
    `, domain.TablePlanningMaterials, domain.TablePurchasedMaterials, domain.TablePlanningMaterialsArchive, domain.TablePurchasedMaterialsArchive, param.SortField, param.Sort)

	rows, err := mr.psql.QueryContext(ctx, sqlQuery, searchPattern, param.CompanyId, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var materials []domain.Material

	for rows.Next() {
		var tableName string //todo подумать надо этим, может возвращать отдельно если потребуется
		var m domain.Material
		err := rows.Scan(&tableName, &m.ID, &m.WarehouseID, &m.Name, &m.ProductCategory, &m.Unit, &m.TotalQuantity, &m.Status, &m.CompanyID)
		if err != nil {
			return nil, 0, err
		}

		materials = append(materials, m)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return materials, totalCount, nil
}

func (mr *MaterialsPostgresRepository) GetIncomeHistoryByWarehouseId(ctx context.Context, id int64, params domain.Param) ([]domain.Material, int64, error) {
	var totalCount int64

	countQuery := fmt.Sprintf(`
	SELECT COUNT(*)
	FROM %s
	WHERE company_id = $1 AND warehouse_id = $2
	`, domain.TablePurchasedMaterials)

	err := mr.psql.QueryRowContext(ctx, countQuery, params.CompanyId, id).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
	SELECT 
	    id, warehouse_id, item_id, name, by_invoice, article, product_category, unit, total_quantity, volume,
		price_without_vat, total_without_vat, supplier_id, location, contract_date, file, status, comments, reserve,
		received_date, last_updated, min_stock_level, expiration_date, responsible_person, storage_cost,
		warehouse_section, incoming_delivery_number, other_fields, company_id, internal_name, units_per_package, 
		supplier_name, contract_number
	FROM %s WHERE company_id = $1 AND warehouse_id = $2 ORDER BY received_date %s LIMIT $3 OFFSET $4
	`, domain.TablePurchasedMaterials, params.Sort)

	rows, err := mr.psql.QueryContext(ctx, query, params.CompanyId, id, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			return
		}
	}(rows)

	var materials []domain.Material

	for rows.Next() {
		var material domain.Material
		var otherFieldsJSON []byte

		if err = rows.Scan(
			&material.ID, &material.WarehouseID, &material.ItemID, &material.Name, &material.ByInvoice, &material.Article,
			pq.Array(&material.ProductCategory), &material.Unit, &material.TotalQuantity, &material.Volume,
			&material.PriceWithoutVAT, &material.TotalWithoutVAT, &material.SupplierID, &material.Location,
			&material.ContractDate, &material.File, &material.Status, &material.Comments, &material.Reserve,
			&material.ReceivedDate, &material.LastUpdated, &material.MinStockLevel, &material.ExpirationDate,
			&material.ResponsiblePerson, &material.StorageCost, &material.WarehouseSection,
			&material.IncomingDeliveryNumber, &otherFieldsJSON, &material.CompanyID, &material.InternalName,
			&material.UnitsPerPackage, &material.SupplierName, &material.ContractNumber,
		); err != nil {
			return nil, 0, err
		}

		if err = json.Unmarshal(otherFieldsJSON, &material.OtherFields); err != nil {
			return nil, 0, err
		}

		materials = append(materials, material)
	}

	return materials, totalCount, nil
}
