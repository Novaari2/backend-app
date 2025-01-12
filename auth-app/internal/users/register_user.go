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

type LoginUserRequest struct {
	Nik      string `json:"nik"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	ID    int    `json:"id"`
	Nik   string `json:"nik"`
	Role  string `json:"role"`
	Token string `json:"token"`
}
