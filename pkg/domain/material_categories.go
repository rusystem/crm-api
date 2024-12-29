package domain

import (
	"time"
)

type MaterialCategory struct {
	ID          int64     `json:"id" example:"1"`
	Name        string    `json:"name" example:"Стальные трубы" binding:"required"`
	CompanyID   int64     `json:"company_id" example:"123"`
	Description string    `json:"description" example:"Категория для стальных труб различного диаметра"`
	Slug        string    `json:"slug" example:"stalnye-truby"`
	CreatedAt   time.Time `json:"created_at" example:"2023-10-01T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-10-02T15:30:00Z"`
	IsActive    bool      `json:"is_active" example:"true"`
	ImgURL      string    `json:"img_url" example:"https://example.com/images/stalnye-truby.jpg"`
}

type CreateMaterialCategory struct {
	Name        string `json:"name" example:"Стальные трубы"`
	Description string `json:"description" example:"Категория для стальных труб различного диаметра"`
	Slug        string `json:"slug" example:"stalnye-truby"`
	ImgURL      string `json:"img_url" example:"https://example.com/images/stalnye-truby.jpg"`
}

type UpdateMaterialCategory struct {
	ID          int64   `json:"-"`
	CompanyID   int64   `json:"-"`
	Name        *string `json:"name" example:"Стальные трубы"`
	Description *string `json:"description" example:"Категория для стальных труб различного диаметра"`
	Slug        *string `json:"slug" example:"stalnye-truby"`
	IsActive    *bool   `json:"is_active" example:"true"`
	ImgURL      *string `json:"img_url" example:"https://example.com/images/stalnye-truby.jpg"`
}
