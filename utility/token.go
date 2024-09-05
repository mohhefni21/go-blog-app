package utility

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(publicId uuid.UUID, role string, secret string, durationInMinute int) (jwtToken string, err error) {

	duration := time.Duration(durationInMinute) * time.Minute
	claims := jwt.MapClaims{
		"publicId": publicId.String(),
		"role":     role,
		"exp":      jwt.NewNumericDate(time.Now().Add(duration)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken, err = token.SignedString([]byte(secret))
	if err != nil {
		return
	}

	return
}

func ValidateToken(tokenString string, secret string) (publicId string, role string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		publicId = fmt.Sprintf("%v", claims["publicId"])
		role = fmt.Sprintf("%v", claims["role"])
		return
	}

	err = fmt.Errorf("unable to extract claims")
	return
}
