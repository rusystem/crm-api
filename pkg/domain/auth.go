package domain

type SignIn struct {
	Username  string `json:"username" binding:"required,min=5,max=140" example:"admin"`
	Password  string `json:"password" binding:"required,min=8,max=255" example:"admin"`
	CompanyId int64  `json:"company_id" example:"1"`
}

type EmailInput struct {
	Email string `json:"email" binding:"required,email,min=5,max=140" example:"example@example.com"`
}

type SignUp struct {
	CompanyId                int64    `json:"company_id" example:"1"`
	Username                 string   `json:"username" binding:"required,min=5,max=140" example:"dmitry"`
	Name                     string   `json:"name" binding:"required,max=140" example:"Дмитрий"`
	Email                    string   `json:"email" binding:"required,email,min=5,max=140" example:"dmitry@test.com"`
	Phone                    string   `json:"phone" binding:"required,min=7,max=140" example:"+77777777777"`
	Password                 string   `json:"password" binding:"required,min=8,max=255" example:"12345678"`
	Role                     string   `json:"role" example:"user"`
	Language                 string   `json:"language" example:"ru"`
	Country                  string   `json:"country" example:"KZ"`
	IsActive                 bool     `json:"is_active" example:"true"`
	IsApproved               bool     `json:"is_approved" example:"true"`
	IsSendSystemNotification bool     `json:"is_send_system_notification" example:"false"`
	Sections                 []string `json:"sections" example:"full_company_access,full_access"`
	Position                 string   `json:"position" example:"test"`
}

type TokensRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
