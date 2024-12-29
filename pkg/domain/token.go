package domain

import "time"

type JWTInfo struct {
	UserId      int64
	CompanyId   int64
	Role        string
	Fingerprint string
	Sections    []string
}

type RefreshSession struct {
	ID        int       `msgpack:"id"`
	UserID    int64     `msgpack:"user_id"`
	CompanyID int64     `msgpack:"company_id"`
	Role      string    `msgpack:"role"`
	Token     string    `msgpack:"token"`
	ExpiresAt time.Time `msgpack:"expires_at"`
	Ip        string    `msgpack:"ip"`
}

type TokenResponse struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
