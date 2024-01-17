package authentication

import (
	"github.com/agung96tm/go-phone-test/internal/models"
	"github.com/pascaldekloe/jwt"
	"os"
	"strconv"
	"time"
)

type JWT struct {
	Claims    jwt.Claims
	SecretKey string
}

func NewJWT(secretKey string) *JWT {
	var claims jwt.Claims
	//claims.Subject = strconv.FormatInt(int64(user.ID), 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))

	return &JWT{
		Claims:    claims,
		SecretKey: secretKey,
	}
}

func (j JWT) GenerateJWT(user *models.User) (string, error) {
	j.Claims.Subject = strconv.FormatInt(int64(user.ID), 10)
	jwtBytes, err := j.Claims.HMACSign(jwt.HS256, []byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return string(jwtBytes), nil
}

func (j JWT) GetClaims(token string) (*jwt.Claims, error) {
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
