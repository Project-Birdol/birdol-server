package jsonmodel

type Auth struct {
	UserID      uint   `json:"user_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
}

type AuthLoginRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
	DeviceID  string `json:"device_id" binding:"required"`
}

type AuthLogoutRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
}
