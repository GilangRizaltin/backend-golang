package pkg

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

type Claims struct {
	Id   int    `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewPayload(id int, role string) *Claims {
	return &Claims{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 40)),
		},
	}
}

func (c *Claims) GenerateToken() (string, error) {
	jwtSecret := os.Getenv("JWT_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	result, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return result, err
}

func VerifyToken(token string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_KEY")
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	payload := parsedToken.Claims.(*Claims)
	return payload, nil
}
