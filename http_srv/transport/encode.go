package transport

type (
	CreateUserResponse struct {
		Id string `json:"Id"`
	}

	AuthenticateResponse struct {
		Token string `json:"Token"`
	}
)
