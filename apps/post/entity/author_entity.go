package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AuthorEntity struct {
	UserId    int            `db:"user_id"`
	PublicId  uuid.UUID      `db:"public_id"`
	Username  string         `db:"username"`
	Fullname  string         `db:"fullname"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	Bio       sql.NullString `db:"bio"`
	Picture   sql.NullString `db:"picture"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
