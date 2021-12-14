package transport

type (
	CreateUserResponse struct {
		Id        string `json:"user_id"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Age       int    `json:"age"`
		ExtraInfo string `json:"extra_information"`
	}

	AuthenticateResponse struct {
		Token string `json:"token"`
	}

	UpdateUserResponse struct {
		Id        int    `json:"user_id"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Age       int    `json:"age"`
		ExtraInfo string `json:"extra_information"`
	}
)
