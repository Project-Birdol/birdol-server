package jsonmodel

type Auth struct {
	UserID      uint   `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type AuthLoginRequest struct {
	Auth     Auth   `json:"auth"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLogoutRequest struct {
	Auth Auth `json:"auth"`
}
