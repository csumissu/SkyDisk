package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.StandardClaims
}

const TokenType = "Bearer"
const ISSUER = "csumissu.xyz"

func NewClaims(userID uint, expirationDuration time.Duration) Claims {
	currentTime := time.Now()
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			Issuer:    ISSUER,
			IssuedAt:  currentTime.Unix(),
			NotBefore: currentTime.Unix(),
			Subject:   strconv.FormatUint(uint64(userID), 10),
			ExpiresAt: currentTime.Add(expirationDuration).Unix(),
		},
	}
	return claims
}

func GenerateJwtToken(signingKey string, claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func ParseJwtToken(signingKey string, token string) (*Claims, error) {
	pureToken := strings.TrimPrefix(token, TokenType+" ")

	jwtToken, err := jwt.Parse(pureToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return convertMapClaimsToLocalClaims(claims), nil
	}

	return nil, fmt.Errorf("jwt token is invalid")
}

func (c Claims) Valid() error {
	var vErr *jwt.ValidationError
	if err := c.StandardClaims.Valid(); err == nil {
		vErr = new(jwt.ValidationError)
	} else {
		vErr = err.(*jwt.ValidationError)
	}

	if !c.VerifyIssuer(ISSUER, true) {
		vErr.Inner = fmt.Errorf("token has no issuer")
		vErr.Errors |= jwt.ValidationErrorIssuer
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

func (c *Claims) GetUserID() uint {
	subject, err := strconv.ParseUint(c.Subject, 10, 32)
	if err != nil {
		Logger.Panic("user id in jwt token is incorrect, subject: %s", c.Subject)
	}
	return uint(subject)
}

func convertMapClaimsToLocalClaims(claims jwt.MapClaims) *Claims {
	return &Claims{
		StandardClaims: jwt.StandardClaims{
			Audience:  GetOrDefault(claims, "aud", "").(string),
			ExpiresAt: int64(claims["exp"].(float64)),
			Id:        claims["jti"].(string),
			IssuedAt:  int64(claims["iat"].(float64)),
			Issuer:    claims["iss"].(string),
			NotBefore: int64(claims["nbf"].(float64)),
			Subject:   claims["sub"].(string),
		},
	}
}
