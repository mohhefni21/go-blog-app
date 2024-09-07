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

type OauthGoogleRequestPayload struct {
	Code  string `json:"code" query:"code"`
	State string `json:"state" query:"state"`
}

type UpdateProfileOnboardingRequestPayload struct {
	Username string `json:"username"`
	Picture  string `json:"picture"`
	Bio      string `json:"bion"`
}
