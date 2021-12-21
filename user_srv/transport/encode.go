package transport

type (
	CreateUserResponse struct {
		Id int
	}

	AuthenticateResponse struct {
		Id int
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

	DeleteUserResponse struct {
		Success bool
	}
)
