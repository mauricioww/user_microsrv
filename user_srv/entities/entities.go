package entities

type (
	User struct {
		Email     string
		Password  string
		Age       int
		ExtraInfo string
	}

	Session struct {
		Email    string
		Password string
		Id       int
	}

	Update struct {
		UserId      int
		Information map[string]interface{}
	}
)
