package jsonmodel

type Auth struct {
	UserID      uint   `json:"user_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
}

type DataLinkRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
	DeviceID  string `json:"device_id" binding:"required"`
	PublicKey string `json:"public_key" binding:"required"`
	KeyType	  string `json:"key_type" binding:"required"`
}

type AuthLogoutRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
}
