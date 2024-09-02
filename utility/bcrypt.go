package utility

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	hashPassowrd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPassowrd), nil
}

// func VerifyPasswordFormPlainn() (err error) {
// 	err = bcrypt.CompareHashAndPassword()
// }
