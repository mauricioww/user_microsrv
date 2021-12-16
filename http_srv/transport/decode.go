package transport

type (
	CreateUserRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Age       int    `json:"age"`
		ExtraInfo string `json:"additional_information"`
	}

	AuthenticateRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UpdateUserRequest struct {
		UserId    int
		Email     string `json:"email"`
		Password  string `json:"password"`
		Age       int    `json:"age"`
		ExtraInfo string `json:"additional_information"`
	}

	GetUserRequest struct {
		UserId int
	}
)
