package domain

// UnitOfMeasure представляет модель единицы измерения
type UnitOfMeasure struct {
	ID           int64  `json:"id" example:"1"`                               // Уникальный идентификатор
	Name         string `json:"name" example:"Килограмм"`                     // Название на языке системы
	NameEn       string `json:"name_en" example:"Kilogram"`                   // Название на английском
	Abbreviation string `json:"abbreviation" example:"kg"`                    // Аббревиатура
	Description  string `json:"description" example:"Единица измерения веса"` // Описание
	CompanyID    int64  `json:"company_id" example:"1"`                       // ID компании
}

type CreateUnitOfMeasure struct {
	Name         string `json:"name" binding:"required,min=1,max=140" example:"Килограмм"`   // Название на языке системы
	NameEn       string `json:"name_en" binding:"required,min=1,max=140" example:"Kilogram"` // Название на английском
	Abbreviation string `json:"abbreviation" example:"kg" binding:"required,min=1,max=140"`  // Аббревиатура
	Description  string `json:"description" example:"Единица измерения веса"`                // Описание
}

type UpdateUnitOfMeasure struct {
	ID           int64   `json:"-"`
	CompanyID    int64   `json:"-"`
	Name         *string `json:"name" example:"Килограмм"`                     // Название на языке системы
	NameEn       *string `json:"name_en" example:"Kilogram"`                   // Название на английском
	Abbreviation *string `json:"abbreviation" example:"kg"`                    // Аббревиатура
	Description  *string `json:"description" example:"Единица измерения веса"` // Описание
}
