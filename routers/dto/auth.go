package dto

import "time"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID      uint                `json:"userId"`
	Nickname    string              `json:"nickname"`
	AccessToken AccessTokenResponse `json:"accessToken"`
}

type AccessTokenResponse struct {
	Type     string    `json:"type"`
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expireAt"`
}
