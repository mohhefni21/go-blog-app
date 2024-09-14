package entity

import (
	"github.com/google/uuid"
)

type UserEntity struct {
	UserId   int       `db:"user_id"`
	PublicId uuid.UUID `db:"public_id"`
	Username string    `db:"username"`
	Email    string    `db:"email"`
}
