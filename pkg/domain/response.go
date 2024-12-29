package domain

type IdResponse struct {
	ID interface{} `json:"id"`
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	IsError bool        `json:"is_error"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type SuccessResponse struct {
	IsError    bool        `json:"is_error"`
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
}

type AvatarResponse struct {
	Avatar string `json:"avatar"`
}

type OperatorStatusResponse struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

type UserStatusResponse struct {
	Status string `json:"status"`
}

type AllowRegistrationResponse struct {
	AllowRegistration bool `json:"allow_registration"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type SignUpResponse struct {
	ID      int64 `json:"id"`
	IsAdmin bool  `json:"is_admin"`
}

type CreateUserResponse struct {
	ID            string `json:"id"`
	CreatedUserID string `json:"created_user_id"`
}
