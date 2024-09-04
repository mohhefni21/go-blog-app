package entity

import "time"

type AuthEntity struct {
	AuthenticationId      int       `db:"authentication_id"`
	UserId                int       `db:"user_id"`
	RefreshToken          string    `db:"refresh_token"`
	RefreshTokenExpiresAt time.Time `db:"refresh_token_expires_at"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}

func NewFromLoginRequestToAuth(userId int, refreshToken string) AuthEntity {
	return AuthEntity{
		UserId:       userId,
		RefreshToken: refreshToken,
	}
}

func (a *AuthEntity) GetRefreshTokenExpiration(refreshTokenExpiresAtInt int) {
	expirationTime := time.Now().Add(time.Duration(refreshTokenExpiresAtInt) * time.Minute)
	a.RefreshTokenExpiresAt = expirationTime
}
