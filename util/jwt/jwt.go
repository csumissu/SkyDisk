package jwt

import (
	"fmt"
	"github.com/csumissu/SkyDisk/conf"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	jwt.StandardClaims
}

func NewClaims() Claims {
	currentTime := time.Now()
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			IssuedAt:  currentTime.Unix(),
			NotBefore: currentTime.Unix(),
			ExpiresAt: currentTime.Add(time.Duration(conf.JwtCfg.ExpirationHours) * time.Hour).Unix(),
		},
	}
	return claims
}

func GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.JwtCfg.SigningKey))
}

func ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(conf.JwtCfg.SigningKey), nil
	})
}
