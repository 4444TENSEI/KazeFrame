package model

import "time"

// 用户资料接口响应体
type UserProfileResponse struct {
	ID        uint      `json:"uid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Signature string    `json:"signature"`
	Gender    string    `json:"gender"`
	AvatarUrl string    `json:"avatar_url"`
	RoleLevel int       `json:"role_level"`
	Online    bool      `json:"online"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

// 删除接口响应体
type UserDeleteResponse struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	SuccessCount int    `json:"successCount"`
	SuccessValue []any  `json:"successValue"`
	ErrorCount   int    `json:"errorCount"`
	ErrorValue   []any  `json:"errorValue"`
}
