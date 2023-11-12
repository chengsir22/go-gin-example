package util

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/luciferCN22/go-gin-example/pkg/setting"
	"time"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		username,
		password,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "gin-blog",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
