package request

type RegisterRequestPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}

type LoginRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegenerateAccessTokenRequestPayload struct {
	RefreshToken string `json:"refreshToken"`
}

type LogoutRequestPayload struct {
	RefreshToken string `json:"refreshToken"`
}
