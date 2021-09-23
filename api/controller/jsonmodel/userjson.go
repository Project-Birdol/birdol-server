package jsonmodel

type UserRequest struct {
	Name      string `json:"name" binding:"required"`
	AccountID string `json:"account_id" binding:"required"`
}

type EditUserRequest struct {
	Id   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type SignupUserRequest struct {
	Name     	string 	`json:"name" binding:"required"`
	PublicKey 	string 	`json:"public_key" binding:"required"`
	DeviceID 	string 	`json:"device_id" binding:"required"`
}

type EnableLinkRequest struct {
	Password  string `json:"password" binding:"required"`
}