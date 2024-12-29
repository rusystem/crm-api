package domain

import "time"

type Company struct {
	ID         int64     `json:"id"`
	NameRu     string    `json:"name_ru"`
	NameEn     string    `json:"name_en"`
	Country    string    `json:"country"`
	Address    string    `json:"address"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Website    string    `json:"website"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsApproved bool      `json:"is_approved"`
	Timezone   string    `json:"timezone"`
}

type CompanyUpdate struct {
	ID         int64   `json:"-"`
	NameRu     *string `json:"name_ru"`
	NameEn     *string `json:"name_en"`
	Country    *string `json:"country"`
	Address    *string `json:"address"`
	Phone      *string `json:"phone"`
	Email      *string `json:"email"`
	Website    *string `json:"website"`
	IsActive   *bool   `json:"is_active"`
	IsApproved *bool   `json:"is_approved"`
	Timezone   *string `json:"timezone"`
}

type CreateCompany struct {
	NameRu     string `json:"name_ru" binding:"required" example:"ООО Рога и копыта"`
	NameEn     string `json:"name_en" example:"OOO ROGA I COPUTA"`
	Country    string `json:"country" example:"KZ"`
	Address    string `json:"address" binding:"required" example:"г. Алматы"`
	Phone      string `json:"phone" binding:"required" example:"+77777777777"`
	Email      string `json:"email" binding:"required" example:"example@example.com"`
	Website    string `json:"website" example:"www.rogakopyta.kz"`
	IsActive   bool   `json:"is_active" example:"true"`
	IsApproved bool   `json:"is_approved" example:"true"`
	Timezone   string `json:"timezone" example:"Asia/Almaty"`
}
