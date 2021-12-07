package transport

type (
	Error struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	CreateUserResponse struct {
		Id        string `json:"id"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Age       int    `json:"Age"`
		ExtraInfo string `json:"extra_information"`
	}

	AuthenticateResponse struct {
		Token string `json:"Token"`
	}
)
