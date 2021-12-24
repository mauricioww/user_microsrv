package entities

type (
	Details struct {
		Country      string  `json:"country"`
		City         string  `json:"city"`
		MobileNumber string  `json:"mobile_number"`
		Married      bool    `json:"married"`
		Height       float32 `json:"height_m"`
		Weigth       float32 `json:"weigth_kg"`
	}

	User struct {
		Email    string
		Password string
		Age      int
		Details
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
