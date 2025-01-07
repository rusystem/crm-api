package domain

import "time"

type Supplier struct {
	ID                int64                  `gorm:"primaryKey" json:"id"`                  // Уникальный идентификатор поставщика
	Name              string                 `json:"name" binding:"required,min=5,max=140"` // Наименование поставщика
	LegalAddress      string                 `json:"legal_address"`                         // Юридический адрес поставщика
	ActualAddress     string                 `json:"actual_address"`                        // Фактический адрес поставщика
	WarehouseAddress  string                 `json:"warehouse_address"`                     // Адрес склада поставщика
	ContactPerson     string                 `json:"contact_person"`                        // Контактное лицо у поставщика
	Phone             string                 `json:"phone" binding:"required"`              // Телефон поставщика
	Email             string                 `json:"email"`                                 // Электронная почта поставщика
	Website           string                 `json:"website"`                               // Сайт поставщика
	ContractNumber    string                 `json:"contract_number"`                       // Номер договора с поставщиком
	ContractDate      time.Time              `json:"contract_date"`                         // Дата договора с поставщиком
	ProductCategories string                 `json:"product_categories"`                    // Категории товаров, поставляемых поставщиком
	PurchaseAmount    float64                `json:"purchase_amount"`                       // Общая сумма закупок у поставщика
	Balance           float64                `json:"balance"`                               // Баланс по поставщику
	ProductTypes      int64                  `json:"product_types"`                         // Количество типов товаров от поставщика
	Comments          string                 `json:"comments"`                              // Комментарии
	Files             string                 `json:"files"`                                 // Ссылки на файлы или документы
	Country           string                 `json:"country"`                               // Страна поставщика
	Region            string                 `json:"region"`                                // Регион или штат поставщика
	TaxID             string                 `json:"tax_id"`                                // Идентификационный номер налогоплательщика (ИНН)
	BankDetails       string                 `json:"bank_details"`                          // Банковские реквизиты поставщика
	RegistrationDate  time.Time              `json:"registration_date"`                     // Дата регистрации поставщика
	PaymentTerms      string                 `json:"payment_terms"`                         // Условия оплаты по контракту
	IsActive          bool                   `json:"is_active"`                             // Статус активности поставщика (активен/неактивен)
	OtherFields       map[string]interface{} `json:"other_fields"`                          // Дополнительные пользовательские поля
	CompanyId         int64                  `json:"company_id"`                            // Уникальный идентификатор компании
}

type InputSupplier struct {
	Name              string                 `json:"name" binding:"required,min=5,max=140" example:"ООО Название поставщика"`          // Наименование поставщика
	LegalAddress      string                 `json:"legal_address" binding:"required" example:"Юридический адрес поставщика"`          // Юридический адрес поставщика
	ActualAddress     string                 `json:"actual_address" binding:"required" example:"Фактический адрес поставщика"`         // Фактический адрес поставщика
	WarehouseAddress  string                 `json:"warehouse_address" binding:"required" example:"Адрес склада поставщика"`           // Адрес склада поставщика
	ContactPerson     string                 `json:"contact_person" binding:"required,max=140" example:"Контактное лицо у поставщика"` // Контактное лицо у поставщика
	Phone             string                 `json:"phone" binding:"required,min=5" example:"Телефон поставщика"`                      // Телефон поставщика
	Email             string                 `json:"email" example:"Электронная почта поставщика"`                                     // Электронная почта поставщика
	Website           string                 `json:"website" example:"Сайт поставщика"`                                                // Сайт поставщика
	ContractNumber    string                 `json:"contract_number" binding:"required" example:"Номер и дата договора с поставщиком"` // Номер и дата договора с поставщиком
	ContractDate      time.Time              `json:"contract_date" binding:"required" example:"2022-01-01T00:00:00Z"`                  // Дата договора с поставщиком
	ProductCategories string                 `json:"product_categories" example:"Категории товаров, поставляемых поставщиком"`         // Категории товаров, поставляемых поставщиком
	PurchaseAmount    float64                `json:"purchase_amount" example:"100.0"`                                                  // Общая сумма закупок у поставщика
	Balance           float64                `json:"balance" binding:"required" example:"100.0"`                                       // Баланс по поставщику
	ProductTypes      int64                  `json:"product_types" example:"100"`                                                      // Количество типов товаров от поставщика
	Comments          string                 `json:"comments" example:"Комментарии"`                                                   // Комментарии
	Files             string                 `json:"files" example:"Ссылки на файлы или документы"`                                    // Ссылки на файлы или документы
	Country           string                 `json:"country" example:"Страна поставщика"`                                              // Страна поставщика
	Region            string                 `json:"region" example:"Регион или штат поставщика"`                                      // Регион или штат поставщика
	TaxID             string                 `json:"tax_id" example:"Идентификационный номер налогоплательщика (ИНН)"`                 // Идентификационный номер налогоплательщика (ИНН)
	BankDetails       string                 `json:"bank_details" binding:"required" example:"Банковские реквизиты поставщика"`        // Банковские реквизиты поставщика
	RegistrationDate  time.Time              `json:"registration_date" example:"2022-01-01T00:00:00Z"`                                 // Дата регистрации поставщика
	PaymentTerms      string                 `json:"payment_terms" example:"Условия оплаты по контракту"`                              // Условия оплаты по контракту
	IsActive          bool                   `json:"is_active" example:"true"`                                                         // Статус активности поставщика (активен/неактивен)
	OtherFields       map[string]interface{} `json:"other_fields"`                                                                     // Дополнительные пользовательские поля
	CompanyId         int64                  `json:"-"`                                                                                // ID компании
}

type UpdateSupplier struct {
	Id                int64                   `json:"-"`                                                                                // ID поставщика
	Name              *string                 `json:"name" binding:"required,min=5,max=140" example:"ООО Название поставщика"`          // Наименование поставщика
	LegalAddress      *string                 `json:"legal_address" binding:"required" example:"Юридический адрес поставщика"`          // Юридический адрес поставщика
	ActualAddress     *string                 `json:"actual_address" binding:"required" example:"Фактический адрес поставщика"`         // Фактический адрес поставщика
	WarehouseAddress  *string                 `json:"warehouse_address" binding:"required" example:"Адрес склада поставщика"`           // Адрес склада поставщика
	ContactPerson     *string                 `json:"contact_person" binding:"required,max=140" example:"Контактное лицо у поставщика"` // Контактное лицо у поставщика
	Phone             *string                 `json:"phone" binding:"required,min=5" example:"Телефон поставщика"`                      // Телефон поставщика
	Email             *string                 `json:"email" example:"Электронная почта поставщика"`                                     // Электронная почта поставщика
	Website           *string                 `json:"website" example:"Сайт поставщика"`                                                // Сайт поставщика
	ContractNumber    *string                 `json:"contract_number" binding:"required" example:"Номер и дата договора с поставщиком"` // Номер и дата договора с поставщиком
	ContractDate      *time.Time              `json:"contract_date" binding:"required" example:"2022-01-01T00:00:00Z"`                  // Дата договора с поставщиком
	ProductCategories *string                 `json:"product_categories" example:"Категории товаров, поставляемых поставщиком"`         // Категории товаров, поставляемых поставщиком
	PurchaseAmount    *float64                `json:"purchase_amount" example:"100.0"`                                                  // Общая сумма закупок у поставщика
	Balance           *float64                `json:"balance" binding:"required" example:"100.0"`                                       // Баланс по поставщику
	ProductTypes      *int64                  `json:"product_types" example:"100"`                                                      // Количество типов товаров от поставщика
	Comments          *string                 `json:"comments" example:"Комментарии"`                                                   // Комментарии
	Files             *string                 `json:"files" example:"Ссылки на файлы или документы"`                                    // Ссылки на файлы или документы
	Country           *string                 `json:"country" example:"Страна поставщика"`                                              // Страна поставщика
	Region            *string                 `json:"region" example:"Регион или штат поставщика"`                                      // Регион или штат поставщика
	TaxID             *string                 `json:"tax_id" example:"Идентификационный номер налогоплательщика (ИНН)"`                 // Идентификационный номер налогоплательщика (ИНН)
	BankDetails       *string                 `json:"bank_details" binding:"required" example:"Банковские реквизиты поставщика"`        // Банковские реквизиты поставщика
	RegistrationDate  *time.Time              `json:"registration_date" example:"2022-01-01T00:00:00Z"`                                 // Дата регистрации поставщика
	PaymentTerms      *string                 `json:"payment_terms" example:"Условия оплаты по контракту"`                              // Условия оплаты по контракту
	IsActive          *bool                   `json:"is_active" example:"true"`                                                         // Статус активности поставщика (активен/неактивен)
	OtherFields       *map[string]interface{} `json:"other_fields"`                                                                     // Дополнительные пользовательские поля
	CompanyId         int64                   `json:"-"`                                                                                // ID компании
}
