package authentication

import (
	"fmt"
	"github.com/agung96tm/go-phone-test/internal/models"
	"github.com/pascaldekloe/jwt"
	"os"
	"strconv"
	"time"
)

type JWT struct {
	claims jwt.Claims
}

func GenerateJWT(user *models.User) (string, error) {
	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(int64(user.ID), 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "foobar"
	}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(secretKey))

	fmt.Println("compare1", secretKey)
	fmt.Println("token1", string(jwtBytes))

	if err != nil {
		return "", err
	}
	return string(jwtBytes), nil
}

func GetClaims(token string) (*jwt.Claims, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "foobar"
	}

	claims, err := jwt.HMACCheck([]byte(token), []byte(secretKey))
	if err != nil {
		return nil, err
	}
	return claims, nil
}
