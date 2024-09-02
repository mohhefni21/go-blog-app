package entity

import (
	"mohhefni/go-blog-app/apps/auth/request"
	"mohhefni/go-blog-app/infra/errorpkg"
	"regexp"
	"time"
)

type Role string

var (
	ROLE_ADMIN Role = "admin"
	ROLE_USER  Role = "user"
)

type AuthEntity struct {
	UserId    int       `db:"userId"`
	Username  string    `db:"username"`
	Fullname  string    `db:"fullname"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      Role      `db:"role"`
	Bio       string    `db:"userId"`
	Picture   string    `db:"picture"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewFromRegisterRequest(req request.RegisterRequestPayload) AuthEntity {
	return AuthEntity{
		Username:  req.Username,
		Fullname:  req.Fullname,
		Email:     req.Email,
		Password:  req.Password,
		Role:      ROLE_USER,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (a *AuthEntity) RegisterValidate() (err error) {
	err = a.UsernameValidate()
	if err != nil {
		return
	}
	err = a.FullnameValidate()
	if err != nil {
		return
	}
	err = a.EmailValidate()
	if err != nil {
		return
	}
	err = a.PasswordValidate()
	if err != nil {
		return
	}

	return
}

func (a *AuthEntity) UsernameValidate() (err error) {
	if a.Username == "" {
		return errorpkg.ErrUsernameRequired
	}

	if len(a.Username) < 3 {
		return errorpkg.ErrUsernameInvalid
	}

	return
}

func (a *AuthEntity) FullnameValidate() (err error) {
	if a.Fullname == "" {
		return errorpkg.ErrFullnameRequired
	}

	return
}

func (a *AuthEntity) EmailValidate() (err error) {
	if a.Email == "" {
		return errorpkg.ErrEmailRequired
	}

	// validate emai with regex
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(a.Email) {
		return errorpkg.ErrEmailInvalid
	}

	return
}

func (a *AuthEntity) PasswordValidate() (err error) {
	if a.Password == "" {
		return errorpkg.ErrPasswordRequired
	}

	if len(a.Username) < 8 {
		return errorpkg.ErrPasswordInvalid
	}

	return
}
