package entity

type AuthEntity struct {
	AuthenticationId      int `db:"authentication_id"`
	UserId                int `db:"user_id"`
	AccessToken           int `db:"access_token"`
	RefreshToken          int `db:"refresh_token"`
	AccessTokenExpiresAt  int `db:"access_token_expires_at"`
	RefreshTokenExpiresAt int `db:"refresh_token_expires_at"`
	CreatedAt             int `db:"created_at"`
	UpdatedAt             int `db:"updated_at"`
}
