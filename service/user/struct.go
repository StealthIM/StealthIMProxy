package user

type registerRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Nickname    string `json:"nickname" validate:"required"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type changePasswordRequest struct {
	Password string `json:"password" validate:"required"`
}

type changeNicknameRequest struct {
	Nickname string `json:"nickname" validate:"required"`
}

type changeEmailRequest struct {
	Email string `json:"email" validate:"required"`
}

type changePhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}
