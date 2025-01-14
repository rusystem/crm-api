package domain

import "time"

// Material представляет структуру товара
type Material struct {
	ID                     int64                  `json:"id"`                       // Уникальный идентификатор записи
	WarehouseID            int64                  `json:"warehouse_id"`             // Склад(место хранения) id
	ItemID                 int64                  `json:"item_id"`                  // Идентификатор товара
	Name                   string                 `json:"name"`                     // Наименование материала от поставщика
	ByInvoice              string                 `json:"by_invoice"`               // Номер товарной накладной
	Article                string                 `json:"article"`                  // Артикул материала
	ProductCategory        []string               `json:"product_category"`         // Категории материала
	Unit                   string                 `json:"unit"`                     // Единица измерения
	TotalQuantity          int64                  `json:"total_quantity"`           // Количество материала
	Volume                 int64                  `json:"volume"`                   // Объем товара
	PriceWithoutVAT        float64                `json:"price_without_vat"`        // Цена без НДС
	TotalWithoutVAT        float64                `json:"total_without_vat"`        // Общая стоимость без НДС
	SupplierID             int64                  `json:"supplier_id"`              // Поставщик товара
	Location               string                 `json:"location"`                 // Локация на складе
	ContractDate           time.Time              `json:"contract_date"`            // Дата договора
	File                   string                 `json:"file"`                     // Файл, связанный с товаром
	Status                 string                 `json:"status"`                   // Статус товара
	Comments               string                 `json:"comments"`                 // Комментарии
	Reserve                string                 `json:"reserve"`                  // Количество на резерве
	ReceivedDate           time.Time              `json:"received_date"`            // Дата поступления на склад
	LastUpdated            time.Time              `json:"last_updated"`             // Дата последнего обновления информации о товаре
	MinStockLevel          int64                  `json:"min_stock_level"`          // Минимальный уровень запаса
	ExpirationDate         time.Time              `json:"expiration_date"`          // Срок годности материала
	ResponsiblePerson      string                 `json:"responsible_person"`       // Ответственное лицо за закуп
	StorageCost            float64                `json:"storage_cost"`             // Стоимость хранения единицы материала
	WarehouseSection       string                 `json:"warehouse_section"`        // Секция хранения
	IncomingDeliveryNumber string                 `json:"incoming_delivery_number"` // Входящий номер поставки
	OtherFields            map[string]interface{} `json:"other_fields"`             // Дополнительные пользовательские поля
	CompanyID              int64                  `json:"company_id"`               // Кабинет компании к кому привязан товар
	InternalName           string                 `json:"internal_name"`            // Наименование материала для внутреннего пользования
	UnitsPerPackage        int64                  `json:"units_per_package"`        // Количество в одной упаковке
	SupplierName           string                 `json:"supplier_name"`            // Поставщик
	ContractNumber         string                 `json:"contract_number"`          // Номер договора
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
	WarehouseID            int64                  `json:"warehouse_id" example:"1"`                                                   // Склад(место хранения) id
	Name                   string                 `json:"name" binding:"required" example:"Steel Beam"`                               // Наименование материала от поставщика
	ByInvoice              string                 `json:"by_invoice" example:"INV-987654"`                                            // Номер товарной накладной
	Article                string                 `json:"article" example:"SB-1234"`                                                  // Артикул материала
	ProductCategory        []string               `json:"product_category" example:"Construction,Bricks"`                             // Категории материала
	Unit                   string                 `json:"unit" example:"pcs"`                                                         // Единица измерения
	TotalQuantity          int64                  `json:"total_quantity" example:"500"`                                               // Количество материала
	Volume                 int64                  `json:"volume" example:"25"`                                                        // Объем товара
	PriceWithoutVAT        float64                `json:"price_without_vat" example:"150.75"`                                         // Цена без НДС
	TotalWithoutVAT        float64                `json:"total_without_vat" example:"75375.00"`                                       // Общая стоимость без НДС
	SupplierID             int64                  `json:"supplier_id" example:"1"`                                                    // Поставщик товара
	ContractDate           time.Time              `json:"contract_date" example:"2023-08-15T10:00:00Z"`                               // Дата договора
	File                   string                 `json:"file" example:"contract_1234.pdf"`                                           // Файл, связанный с товаром
	Status                 string                 `json:"status" example:"active"`                                                    // Статус товара
	Comments               string                 `json:"comments" example:"Urgent order"`                                            // Комментарии
	Reserve                string                 `json:"reserve" example:"50"`                                                       // Количество на резерве
	ReceivedDate           time.Time              `json:"received_date" example:"2023-08-20T10:00:00Z"`                               // Дата поступления на склад
	MinStockLevel          int64                  `json:"min_stock_level" example:"10"`                                               // Минимальный уровень запаса
	ExpirationDate         time.Time              `json:"expiration_date" example:"2024-08-15T10:00:00Z"`                             // Срок годности материала
	ResponsiblePerson      string                 `json:"responsible_person" example:"John Doe"`                                      // Ответственное лицо за закуп
	StorageCost            float64                `json:"storage_cost" example:"500.00"`                                              // Стоимость хранения единицы материала
	WarehouseSection       string                 `json:"warehouse_section" example:"B-Section-2"`                                    // Секция хранения
	IncomingDeliveryNumber string                 `json:"incoming_delivery_number" example:"DEL-56789"`                               // Входящий номер поставки
	OtherFields            map[string]interface{} `json:"other_fields"`                                                               // Дополнительные пользовательские поля
	InternalName           string                 `json:"internal_name" example:"Наименование материала для внутреннего пользования"` // Наименование материала для внутреннего пользования
	UnitsPerPackage        int64                  `json:"units_per_package" example:"4"`                                              // Количество в одной упаковке
	SupplierName           string                 `json:"supplier_name" example:"Поставщик"`                                          // Поставщик
	ContractNumber         string                 `json:"contract_number" example:"11"`                                               // Номер договора
}

// UpdatePlanningMaterial представляет структуру товара
type UpdatePlanningMaterial struct {
	ID                     int64                   `json:"-"`                                                                          // Уникальный идентификатор записи
	WarehouseID            *int64                  `json:"warehouse_id" example:"1"`                                                   // Склад(место хранения) id
	Name                   *string                 `json:"name" example:"Steel Beam"`                                                  // Наименование материала от поставщика
	ByInvoice              *string                 `json:"by_invoice" example:"INV-987654"`                                            // Номер товарной накладной
	Article                *string                 `json:"article" example:"SB-1234"`                                                  // Артикул материала
	ProductCategory        *[]string               `json:"product_category" example:"Construction,Bricks"`                             // Категории материала
	Unit                   *string                 `json:"unit" example:"pcs"`                                                         // Единица измерения
	TotalQuantity          *int64                  `json:"total_quantity" example:"500"`                                               // Количество материала
	Volume                 *int64                  `json:"volume" example:"25"`                                                        // Объем товара
	PriceWithoutVAT        *float64                `json:"price_without_vat" example:"150.75"`                                         // Цена без НДС
	TotalWithoutVAT        *float64                `json:"total_without_vat" example:"75375.00"`                                       // Общая стоимость без НДС
	SupplierID             *int64                  `json:"supplier_id" example:"1"`                                                    // Поставщик товара
	Location               *string                 `json:"location" example:"A1-Section-3"`                                            // Локация на складе
	ContractDate           *time.Time              `json:"contract_date" example:"2023-08-15T10:00:00Z"`                               // Дата договора
	File                   *string                 `json:"file" example:"contract_1234.pdf"`                                           // Файл, связанный с товаром
	Status                 *string                 `json:"status" example:"active"`                                                    // Статус товара
	Comments               *string                 `json:"comments" example:"Urgent order"`                                            // Комментарии
	Reserve                *string                 `json:"reserve" example:"50"`                                                       // Количество на резерве
	ReceivedDate           *time.Time              `json:"received_date" example:"2023-08-20T10:00:00Z"`                               // Дата поступления на склад
	MinStockLevel          *int64                  `json:"min_stock_level" example:"10"`                                               // Минимальный уровень запаса
	ExpirationDate         *time.Time              `json:"expiration_date" example:"2024-08-15T10:00:00Z"`                             // Срок годности материала
	ResponsiblePerson      *string                 `json:"responsible_person" example:"John Doe"`                                      // Ответственное лицо за закуп
	StorageCost            *float64                `json:"storage_cost" example:"500.00"`                                              // Стоимость хранения единицы материала
	WarehouseSection       *string                 `json:"warehouse_section" example:"B-Section-2"`                                    // Секция хранения
	IncomingDeliveryNumber *string                 `json:"incoming_delivery_number" example:"DEL-56789"`                               // Входящий номер поставки
	OtherFields            *map[string]interface{} `json:"other_fields"`                                                               // Дополнительные пользовательские поля
	InternalName           *string                 `json:"internal_name" example:"Наименование материала для внутреннего пользования"` // Наименование материала для внутреннего пользования
	UnitsPerPackage        *int64                  `json:"units_per_package" example:"10"`                                             // Количество в одной упаковке
	SupplierName           *string                 `json:"supplier_name" example:"Поставщик 2"`                                        // Поставщик
	ContractNumber         *string                 `json:"contract_number" example:"22"`                                               // Номер договора
}

type PurchasedIdResponse struct {
	ID     int64 `json:"id"`
	ItemId int64 `json:"item_id"`
}

// CreatePurchasedMaterial представляет структуру товара
type CreatePurchasedMaterial struct {
	WarehouseID            int64                  `json:"warehouse_id" example:"1"`                                                   // Склад(место хранения) id
	Name                   string                 `json:"name" binding:"required" example:"Steel Beam"`                               // Наименование материала от поставщика
	ByInvoice              string                 `json:"by_invoice" example:"INV-987654"`                                            // Номер товарной накладной
	Article                string                 `json:"article" example:"SB-1234"`                                                  // Артикул материала
	ProductCategory        []string               `json:"product_category" example:"Construction,Bricks"`                             // Категории материала
	Unit                   string                 `json:"unit" example:"pcs"`                                                         // Единица измерения
	TotalQuantity          int64                  `json:"total_quantity" example:"500"`                                               // Количество материала
	Volume                 int64                  `json:"volume" example:"25"`                                                        // Объем товара
	PriceWithoutVAT        float64                `json:"price_without_vat" example:"150.75"`                                         // Цена без НДС
	TotalWithoutVAT        float64                `json:"total_without_vat" example:"75375.00"`                                       // Общая стоимость без НДС
	SupplierID             int64                  `json:"supplier_id" example:"1"`                                                    // Поставщик товара
	Location               string                 `json:"location" example:"A1-Section-3"`                                            // Локация на складе
	ContractDate           time.Time              `json:"contract_date" example:"2023-08-15T10:00:00Z"`                               // Дата договора
	File                   string                 `json:"file" example:"contract_1234.pdf"`                                           // Файл, связанный с товаром
	Status                 string                 `json:"status" example:"active"`                                                    // Статус товара
	Comments               string                 `json:"comments" example:"Urgent order"`                                            // Комментарии
	Reserve                string                 `json:"reserve" example:"50"`                                                       // Количество на резерве
	ReceivedDate           time.Time              `json:"received_date" example:"2023-08-20T10:00:00Z"`                               // Дата поступления на склад
	MinStockLevel          int64                  `json:"min_stock_level" example:"10"`                                               // Минимальный уровень запаса
	ExpirationDate         time.Time              `json:"expiration_date" example:"2024-08-15T10:00:00Z"`                             // Срок годности материала
	ResponsiblePerson      string                 `json:"responsible_person" example:"John Doe"`                                      // Ответственное лицо за закуп
	StorageCost            float64                `json:"storage_cost" example:"500.00"`                                              // Стоимость хранения единицы материала
	WarehouseSection       string                 `json:"warehouse_section" example:"B-Section-2"`                                    // Секция хранения
	IncomingDeliveryNumber string                 `json:"incoming_delivery_number" example:"DEL-56789"`                               // Входящий номер поставки
	OtherFields            map[string]interface{} `json:"other_fields"`                                                               // Дополнительные пользовательские поля
	InternalName           string                 `json:"internal_name" example:"Наименование материала для внутреннего пользования"` // Наименование материала для внутреннего пользования
	UnitsPerPackage        int64                  `json:"units_per_package" example:"2"`                                              // Количество в одной упаковке
	SupplierName           string                 `json:"supplier_name" example:"Поставщик 3"`                                        // Поставщик
	ContractNumber         string                 `json:"contract_number" example:"33"`                                               // Номер договора
}

// UpdatePurchasedMaterial представляет структуру товара
type UpdatePurchasedMaterial struct {
	ID                     int64                   `json:"-"`                                                                          // Уникальный идентификатор записи
	WarehouseID            *int64                  `json:"warehouse_id" example:"1"`                                                   // Склад(место хранения) id
	Name                   *string                 `json:"name" example:"Steel Beam"`                                                  // Наименование материала от поставщика
	ByInvoice              *string                 `json:"by_invoice" example:"INV-987654"`                                            // Номер товарной накладной
	Article                *string                 `json:"article" example:"SB-1234"`                                                  // Артикул материала
	ProductCategory        *[]string               `json:"product_category" example:"Construction,Bricks"`                             // Категории материала
	Unit                   *string                 `json:"unit" example:"pcs"`                                                         // Единица измерения
	TotalQuantity          *int64                  `json:"total_quantity" example:"500"`                                               // Количество материала
	Volume                 *int64                  `json:"volume" example:"25"`                                                        // Объем товара
	PriceWithoutVAT        *float64                `json:"price_without_vat" example:"150.75"`                                         // Цена без НДС
	TotalWithoutVAT        *float64                `json:"total_without_vat" example:"75375.00"`                                       // Общая стоимость без НДС
	SupplierID             *int64                  `json:"supplier_id" example:"1"`                                                    // Поставщик товара
	Location               *string                 `json:"location" example:"A1-Section-3"`                                            // Локация на складе
	ContractDate           *time.Time              `json:"contract_date" example:"2023-08-15T10:00:00Z"`                               // Дата договора
	File                   *string                 `json:"file" example:"contract_1234.pdf"`                                           // Файл, связанный с товаром
	Status                 *string                 `json:"status" example:"active"`                                                    // Статус товара
	Comments               *string                 `json:"comments" example:"Urgent order"`                                            // Комментарии
	Reserve                *string                 `json:"reserve" example:"50"`                                                       // Количество на резерве
	ReceivedDate           *time.Time              `json:"received_date" example:"2023-08-20T10:00:00Z"`                               // Дата поступления на склад
	MinStockLevel          *int64                  `json:"min_stock_level" example:"10"`                                               // Минимальный уровень запаса
	ExpirationDate         *time.Time              `json:"expiration_date" example:"2024-08-15T10:00:00Z"`                             // Срок годности материала
	ResponsiblePerson      *string                 `json:"responsible_person" example:"John Doe"`                                      // Ответственное лицо за закуп
	StorageCost            *float64                `json:"storage_cost" example:"500.00"`                                              // Стоимость хранения единицы материала
	WarehouseSection       *string                 `json:"warehouse_section" example:"B-Section-2"`                                    // Секция хранения
	IncomingDeliveryNumber *string                 `json:"incoming_delivery_number" example:"DEL-56789"`                               // Входящий номер поставки
	OtherFields            *map[string]interface{} `json:"other_fields"`                                                               // Дополнительные пользовательские поля
	InternalName           *string                 `json:"internal_name" example:"Наименование материала для внутреннего пользования"` // Наименование материала для внутреннего пользования
	UnitsPerPackage        *int64                  `json:"units_per_package" example:"10"`                                             // Количество в одной упаковке
	SupplierName           *string                 `json:"supplier_name" example:"Поставщик 3"`                                        // Поставщик
	ContractNumber         *string                 `json:"contract_number" example:"44"`                                               // Номер договора
}
