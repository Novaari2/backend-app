package users

type RegisterUserRequest struct {
	Nik  string `json:"nik"`
	Role string `json:"role"`
}

type RegisterUserResponse struct {
	Nik      string `json:"nik"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
