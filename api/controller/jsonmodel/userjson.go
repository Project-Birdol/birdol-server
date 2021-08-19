package jsonmodel

type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EditUserRequest struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SignupUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// TODO: Add DeviceID
}
