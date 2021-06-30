package jwt

import (
	"fmt"
	"github.com/csumissu/SkyDisk/conf"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"strings"
	"time"
)

var TokenType = "Bearer"
var ISSUER = "csumissu.xyz"

func DefaultExpirationDuration() time.Duration {
	return time.Duration(conf.JwtCfg.ExpirationHours) * time.Hour
}

func NewClaims(payload map[string]interface{}) jwt.MapClaims {
	currentTime := time.Now()
	claims := make(jwt.MapClaims)
	claims["jti"] = uuid.NewString()
	claims["iat"] = currentTime.Unix()
	claims["nbf"] = currentTime.Unix()
	claims["exp"] = currentTime.Add(DefaultExpirationDuration()).Unix()
	claims["iss"] = ISSUER

	for key, value := range payload {
		claims[key] = value
	}

	return claims
}

func GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.JwtCfg.SigningKey))
}

func ParseToken(token string) (jwt.MapClaims, error) {
	pureToken := strings.ReplaceAll(token, TokenType+" ", "")

	jwtToken, err := jwt.Parse(pureToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(conf.JwtCfg.SigningKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("jwt token is invalid")
}
