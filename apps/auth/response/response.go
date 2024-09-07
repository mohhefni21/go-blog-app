package response

type OauthResponse struct {
	AccessToken  string
	RefreshToken string
	Email        string
}

func NewOauthResponse(accessToken string, refreshToken string, email string) OauthResponse {
	return OauthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        email,
	}
}
