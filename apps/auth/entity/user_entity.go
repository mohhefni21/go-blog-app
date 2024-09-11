package entity

import (
	"database/sql"
	"fmt"
	"mohhefni/go-blog-app/apps/auth/request"
	"mohhefni/go-blog-app/infra/errorpkg"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Role string

var (
	ROLE_ADMIN Role = "admin"
	ROLE_USER  Role = "user"
)

type UserEntity struct {
	UserId    int            `db:"user_id"`
	PublicId  uuid.UUID      `db:"public_id"`
	Username  string         `db:"username"`
	Fullname  string         `db:"fullname"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	Role      Role           `db:"role"`
	Bio       sql.NullString `db:"bio"`
	Picture   sql.NullString `db:"picture"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func NewFromRegisterRequest(req request.RegisterRequestPayload) UserEntity {
	return UserEntity{
		PublicId:  uuid.New(),
		Username:  req.Username,
		Fullname:  req.Fullname,
		Email:     req.Email,
		Password:  req.Password,
		Role:      ROLE_USER,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewFromLoginRequest(req request.LoginRequestPayload) UserEntity {
	return UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (a *UserEntity) GenerateUsernameOauth(id string) {
	toLowerCase := strings.ToLower(a.Username)
	underscoreReplace := strings.ReplaceAll(toLowerCase, " ", "_")
	reg := regexp.MustCompile(`[^a-z0-9_]+`)
	cleanedCaracter := reg.ReplaceAllString(underscoreReplace, "")
	usernameGenerated := fmt.Sprintf("%s_%s", cleanedCaracter, id)
	a.Username = usernameGenerated
}

func (a *UserEntity) ValidateLogin() (err error) {
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

func (a *UserEntity) RegisterValidate() (err error) {
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

func (a *UserEntity) UsernameValidate() (err error) {
	if a.Username == "" {
		return errorpkg.ErrUsernameRequired
	}

	if len(a.Username) < 3 {
		return errorpkg.ErrUsernameInvalid
	}

	// validasi with regex
	re := regexp.MustCompile("^[a-zA-Z0-9]([._]?[a-zA-Z0-9]+)*$")
	matching := re.MatchString(a.Username)
	if !matching {
		return errorpkg.ErrUsernameInvalid
	}

	return
}

func (a *UserEntity) FullnameValidate() (err error) {
	if a.Fullname == "" {
		return errorpkg.ErrFullnameRequired
	}

	return
}

func (a *UserEntity) EmailValidate() (err error) {
	if a.Email == "" {
		return errorpkg.ErrEmailRequired
	}

	// validation email with regex
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(a.Email) {
		return errorpkg.ErrEmailInvalid
	}

	return
}

func (a *UserEntity) PasswordValidate() (err error) {
	if a.Password == "" {
		return errorpkg.ErrPasswordRequired
	}

	if len(a.Password) < 8 {
		return errorpkg.ErrPasswordInvalid
	}

	return
}
