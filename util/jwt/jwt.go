package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/csumissu/SkyDisk/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var TokenType = "Bearer"
var ISSUER = "csumissu.xyz"

func DefaultExpirationDuration() time.Duration {
	return time.Duration(config.JwtCfg.ExpirationHours) * time.Hour
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
	return token.SignedString([]byte(config.JwtCfg.SigningKey))
}

func ParseToken(token string) (jwt.MapClaims, error) {
	pureToken := strings.TrimPrefix(token, TokenType+" ")

	jwtToken, err := jwt.Parse(pureToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(config.JwtCfg.SigningKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("jwt token is invalid")
}
