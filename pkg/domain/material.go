package domain

import "time"

// Material представляет структуру товара
type Material struct {
	ID                     int64                  `json:"id"`                       // Уникальный идентификатор записи
	WarehouseID            int64                  `json:"warehouse_id"`             // Id склада
	ItemID                 int64                  `json:"item_id"`                  // Идентификатор товара
	Name                   string                 `json:"name"`                     // Наименование товара
	ByInvoice              string                 `json:"by_invoice"`               // Накладная на товар
	Article                string                 `json:"article"`                  // Артикул товара
	ProductCategory        string                 `json:"product_category"`         // Категория товара
	Unit                   string                 `json:"unit"`                     // Единица измерения
	TotalQuantity          int64                  `json:"total_quantity"`           // Общее количество товара
	Volume                 int64                  `json:"volume"`                   // Объем товара
	PriceWithoutVAT        float64                `json:"price_without_vat"`        // Цена без НДС
	TotalWithoutVAT        float64                `json:"total_without_vat"`        // Общая стоимость без НДС
	SupplierID             int64                  `json:"supplier_id"`              // Поставщик товара
	Location               string                 `json:"location"`                 // Локация на складе
	Contract               time.Time              `json:"contract"`                 // Дата договора
	File                   string                 `json:"file"`                     // Файл, связанный с товаром
	Status                 string                 `json:"status"`                   // Статус товара
	Comments               string                 `json:"comments"`                 // Комментарии
	Reserve                string                 `json:"reserve"`                  // Резерв товара
	ReceivedDate           time.Time              `json:"received_date"`            // Дата поступления товара
	LastUpdated            time.Time              `json:"last_updated"`             // Дата последнего обновления информации о товаре
	MinStockLevel          int64                  `json:"min_stock_level"`          // Минимальный уровень запаса
	ExpirationDate         time.Time              `json:"expiration_date"`          // Срок годности товара
	ResponsiblePerson      string                 `json:"responsible_person"`       // Ответственное лицо за товар
	StorageCost            float64                `json:"storage_cost"`             // Стоимость хранения товара
	WarehouseSection       string                 `json:"warehouse_section"`        // Секция склада, где хранится товар
	IncomingDeliveryNumber string                 `json:"incoming_delivery_number"` // Входящий номер поставки
	OtherFields            map[string]interface{} `json:"other_fields"`             // Дополнительные пользовательские поля
	CompanyID              int64                  `json:"company_id"`               // Кабинет компании к кому привязан товар
}

type MaterialParams struct {
	Limit     int64
	Offset    int64
	CompanyId int64
	Query     string
	Sort      string `json:"sort"`
	SortField string `json:"sort_field"`
}

// CreatePlanningMaterial представляет структуру создания товара
type CreatePlanningMaterial struct {
	WarehouseID            int64                  `json:"warehouse_id" example:"1"`                       // Id склада
	Name                   string                 `json:"name" binding:"required" example:"Steel Beam"`   // Наименование товара
	ByInvoice              string                 `json:"by_invoice" example:"INV-987654"`                // Накладная на товар
	Article                string                 `json:"article" example:"SB-1234"`                      // Артикул товара
	ProductCategory        string                 `json:"product_category" example:"Construction"`        // Категория товара
	Unit                   string                 `json:"unit" example:"pcs"`                             // Единица измерения
	TotalQuantity          int64                  `json:"total_quantity" example:"500"`                   // Общее количество товара
	Volume                 int64                  `json:"volume" example:"25"`                            // Объем товара
	PriceWithoutVAT        float64                `json:"price_without_vat" example:"150.75"`             // Цена без НДС
	TotalWithoutVAT        float64                `json:"total_without_vat" example:"75375.00"`           // Общая стоимость без НДС
	SupplierID             int64                  `json:"supplier_id" example:"1"`                        // Поставщик товара
	Contract               time.Time              `json:"contract" example:"2023-08-15T10:00:00Z"`        // Дата договора
	File                   string                 `json:"file" example:"contract_1234.pdf"`               // Файл, связанный с товаром
	Status                 string                 `json:"status" example:"active"`                        // Статус товара
	Comments               string                 `json:"comments" example:"Urgent order"`                // Комментарии
	Reserve                string                 `json:"reserve" example:"50"`                           // Резерв товара
	ReceivedDate           time.Time              `json:"received_date" example:"2023-08-20T10:00:00Z"`   // Дата поступления товара
	MinStockLevel          int64                  `json:"min_stock_level" example:"10"`                   // Минимальный уровень запаса
	ExpirationDate         time.Time              `json:"expiration_date" example:"2024-08-15T10:00:00Z"` // Срок годности товара
	ResponsiblePerson      string                 `json:"responsible_person" example:"John Doe"`          // Ответственное лицо за товар
	StorageCost            float64                `json:"storage_cost" example:"500.00"`                  // Стоимость хранения товара
	WarehouseSection       string                 `json:"warehouse_section" example:"B-Section-2"`        // Секция склада, где хранится товар
	IncomingDeliveryNumber string                 `json:"incoming_delivery_number" example:"DEL-56789"`   // Входящий номер поставки
	OtherFields            map[string]interface{} `json:"other_fields"`                                   // Дополнительные пользовательские поля
}

// UpdatePlanningMaterial представляет структуру товара
type UpdatePlanningMaterial struct {
	ID                     int64                   `json:"-"`                                              // Уникальный идентификатор записи
	WarehouseID            *int64                  `json:"warehouse_id" example:"1"`                       // Id склада
	Name                   *string                 `json:"name" example:"Steel Beam"`                      // Наименование товара
	ByInvoice              *string                 `json:"by_invoice" example:"INV-987654"`                // Накладная на товар
	Article                *string                 `json:"article" example:"SB-1234"`                      // Артикул товара
	ProductCategory        *string                 `json:"product_category" example:"Construction"`        // Категория товара
	Unit                   *string                 `json:"unit" example:"pcs"`                             // Единица измерения
	TotalQuantity          *int64                  `json:"total_quantity" example:"500"`                   // Общее количество товара
	Volume                 *int64                  `json:"volume" example:"25"`                            // Объем товара
	PriceWithoutVAT        *float64                `json:"price_without_vat" example:"150.75"`             // Цена без НДС
	TotalWithoutVAT        *float64                `json:"total_without_vat" example:"75375.00"`           // Общая стоимость без НДС
	SupplierID             *int64                  `json:"supplier_id" example:"1"`                        // Поставщик товара
	Location               *string                 `json:"location" example:"A1-Section-3"`                // Локация на складе
	Contract               *time.Time              `json:"contract" example:"2023-08-15T10:00:00Z"`        // Дата договора
	File                   *string                 `json:"file" example:"contract_1234.pdf"`               // Файл, связанный с товаром
	Status                 *string                 `json:"status" example:"active"`                        // Статус товара
	Comments               *string                 `json:"comments" example:"Urgent order"`                // Комментарии
	Reserve                *string                 `json:"reserve" example:"50"`                           // Резерв товара
	ReceivedDate           *time.Time              `json:"received_date" example:"2023-08-20T10:00:00Z"`   // Дата поступления товара
	MinStockLevel          *int64                  `json:"min_stock_level" example:"10"`                   // Минимальный уровень запаса
	ExpirationDate         *time.Time              `json:"expiration_date" example:"2024-08-15T10:00:00Z"` // Срок годности товара
	ResponsiblePerson      *string                 `json:"responsible_person" example:"John Doe"`          // Ответственное лицо за товар
	StorageCost            *float64                `json:"storage_cost" example:"500.00"`                  // Стоимость хранения товара
	WarehouseSection       *string                 `json:"warehouse_section" example:"B-Section-2"`        // Секция склада, где хранится товар
	IncomingDeliveryNumber *string                 `json:"incoming_delivery_number" example:"DEL-56789"`   // Входящий номер поставки
	OtherFields            *map[string]interface{} `json:"other_fields"`                                   // Дополнительные пользовательские поля
}

type PurchasedIdResponse struct {
	ID     int64 `json:"id"`
	ItemId int64 `json:"item_id"`
}

// CreatePurchasedMaterial представляет структуру товара
type CreatePurchasedMaterial struct {
	WarehouseID            int64                  `json:"warehouse_id" example:"1"`                       // Id склада
	Name                   string                 `json:"name" binding:"required" example:"Steel Beam"`   // Наименование товара
	ByInvoice              string                 `json:"by_invoice" example:"INV-987654"`                // Накладная на товар
	Article                string                 `json:"article" example:"SB-1234"`                      // Артикул товара
	ProductCategory        string                 `json:"product_category" example:"Construction"`        // Категория товара
	Unit                   string                 `json:"unit" example:"pcs"`                             // Единица измерения
	TotalQuantity          int64                  `json:"total_quantity" example:"500"`                   // Общее количество товара
	Volume                 int64                  `json:"volume" example:"25"`                            // Объем товара
	PriceWithoutVAT        float64                `json:"price_without_vat" example:"150.75"`             // Цена без НДС
	TotalWithoutVAT        float64                `json:"total_without_vat" example:"75375.00"`           // Общая стоимость без НДС
	SupplierID             int64                  `json:"supplier_id" example:"1"`                        // Поставщик товара
	Location               string                 `json:"location" example:"A1-Section-3"`                // Локация на складе
	Contract               time.Time              `json:"contract" example:"2023-08-15T10:00:00Z"`        // Дата договора
	File                   string                 `json:"file" example:"contract_1234.pdf"`               // Файл, связанный с товаром
	Status                 string                 `json:"status" example:"active"`                        // Статус товара
	Comments               string                 `json:"comments" example:"Urgent order"`                // Комментарии
	Reserve                string                 `json:"reserve" example:"50"`                           // Резерв товара
	ReceivedDate           time.Time              `json:"received_date" example:"2023-08-20T10:00:00Z"`   // Дата поступления товара
	MinStockLevel          int64                  `json:"min_stock_level" example:"10"`                   // Минимальный уровень запаса
	ExpirationDate         time.Time              `json:"expiration_date" example:"2024-08-15T10:00:00Z"` // Срок годности товара
	ResponsiblePerson      string                 `json:"responsible_person" example:"John Doe"`          // Ответственное лицо за товар
	StorageCost            float64                `json:"storage_cost" example:"500.00"`                  // Стоимость хранения товара
	WarehouseSection       string                 `json:"warehouse_section" example:"B-Section-2"`        // Секция склада, где хранится товар
	IncomingDeliveryNumber string                 `json:"incoming_delivery_number" example:"DEL-56789"`   // Входящий номер поставки
	OtherFields            map[string]interface{} `json:"other_fields"`                                   // Дополнительные пользовательские поля
}

// UpdatePurchasedMaterial представляет структуру товара
type UpdatePurchasedMaterial struct {
	ID                     int64                   `json:"-"`                                              // Уникальный идентификатор записи
	WarehouseID            *int64                  `json:"warehouse_id" example:"1"`                       // Id склада
	Name                   *string                 `json:"name" example:"Steel Beam"`                      // Наименование товара
	ByInvoice              *string                 `json:"by_invoice" example:"INV-987654"`                // Накладная на товар
	Article                *string                 `json:"article" example:"SB-1234"`                      // Артикул товара
	ProductCategory        *string                 `json:"product_category" example:"Construction"`        // Категория товара
	Unit                   *string                 `json:"unit" example:"pcs"`                             // Единица измерения
	TotalQuantity          *int64                  `json:"total_quantity" example:"500"`                   // Общее количество товара
	Volume                 *int64                  `json:"volume" example:"25"`                            // Объем товара
	PriceWithoutVAT        *float64                `json:"price_without_vat" example:"150.75"`             // Цена без НДС
	TotalWithoutVAT        *float64                `json:"total_without_vat" example:"75375.00"`           // Общая стоимость без НДС
	SupplierID             *int64                  `json:"supplier_id" example:"1"`                        // Поставщик товара
	Location               *string                 `json:"location" example:"A1-Section-3"`                // Локация на складе
	Contract               *time.Time              `json:"contract" example:"2023-08-15T10:00:00Z"`        // Дата договора
	File                   *string                 `json:"file" example:"contract_1234.pdf"`               // Файл, связанный с товаром
	Status                 *string                 `json:"status" example:"active"`                        // Статус товара
	Comments               *string                 `json:"comments" example:"Urgent order"`                // Комментарии
	Reserve                *string                 `json:"reserve" example:"50"`                           // Резерв товара
	ReceivedDate           *time.Time              `json:"received_date" example:"2023-08-20T10:00:00Z"`   // Дата поступления товара
	MinStockLevel          *int64                  `json:"min_stock_level" example:"10"`                   // Минимальный уровень запаса
	ExpirationDate         *time.Time              `json:"expiration_date" example:"2024-08-15T10:00:00Z"` // Срок годности товара
	ResponsiblePerson      *string                 `json:"responsible_person" example:"John Doe"`          // Ответственное лицо за товар
	StorageCost            *float64                `json:"storage_cost" example:"500.00"`                  // Стоимость хранения товара
	WarehouseSection       *string                 `json:"warehouse_section" example:"B-Section-2"`        // Секция склада, где хранится товар
	IncomingDeliveryNumber *string                 `json:"incoming_delivery_number" example:"DEL-56789"`   // Входящий номер поставки
	OtherFields            *map[string]interface{} `json:"other_fields"`                                   // Дополнительные пользовательские поля
}
