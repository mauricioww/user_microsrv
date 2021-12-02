package transport

type (
	UserRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Age       int    `json:"age"`
		ExtraInfo string `json:"additional_information"`
	}
)
