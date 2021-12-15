package transport

type (
	CreateUserResponse struct {
		Id string
	}

	AuthenticateResponse struct {
		Id string
	}

	UpdateUserResponse struct {
		Email     string
		Password  string
		Age       int
		ExtraInfo string
	}

	GetUserResponse struct {
		Email     string
		Password  string
		Age       int
		ExtraInfo string
	}
)
