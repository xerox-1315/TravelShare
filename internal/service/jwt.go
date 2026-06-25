package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xerox-1315/TravelShare.git/errs"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", errs.ErrorGenerateToken
	}
	return signed, nil
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, errs.ErrorInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errs.ErrorInvalidToken
	}
	return claims, nil
}
