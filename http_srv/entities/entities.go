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
	}

	UserUpdate struct {
		UserId int
		User
	}
)
