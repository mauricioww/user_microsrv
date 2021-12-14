package transport

type (
	CreateUserRequest struct {
		Email     string
		Password  string
		Age       int
		ExtraInfo string
	}

	AuthenticateRequest struct {
		Email    string
		Password string
	}

	UpdateUserRequest struct {
		Id        int
		Email     string
		Password  string
		Age       int
		ExtraInfo string
	}
)
