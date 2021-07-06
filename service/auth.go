package service

import (
	"fmt"
	"github.com/csumissu/SkyDisk/config"
	"github.com/csumissu/SkyDisk/infra"
	"github.com/csumissu/SkyDisk/models"
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/util"
	"net/http"
	"time"
)

type AuthService struct {
}

func (service *AuthService) Login(input dto.LoginRequest) dto.Response {
	user, err := models.GetUserByUsername(input.Username)
	if err != nil {
		return dto.Failure(http.StatusBadRequest, "username or password is incorrect")
	}
	if authOK, _ := user.CheckPassword(input.Password); !authOK {
		return dto.Failure(http.StatusBadRequest, "username or password is incorrect")
	}
	if user.Status == models.Banned {
		return dto.Failure(http.StatusForbidden, "this user is banned")
	}

	expirationDuration := time.Duration(config.JwtCfg.ExpirationHours) * time.Hour
	claims := util.NewClaims(user.ID, expirationDuration)
	token, err := util.GenerateJwtToken(config.JwtCfg.SigningKey, claims)
	if err != nil {
		return dto.Failure(http.StatusInternalServerError, "token could not be generated")
	}
	infra.RedisClient.Set(claims.Id, user.ID, expirationDuration)

	response := &dto.LoginResponse{
		UserID:   int(user.ID),
		Nickname: user.Nickname,
		AccessToken: dto.AccessTokenResponse{
			Type:     util.TokenType,
			Token:    token,
			ExpireAt: time.Unix(claims.ExpiresAt, 0),
		},
	}
	return dto.Success(response)
}

func (service *AuthService) Logout(token string) dto.Response {
	if claims, err := CheckAuthorizationHeader(token); err != nil {
		return dto.Failure(http.StatusUnauthorized, err.Error())
	} else {
		infra.RedisClient.Del(claims.Id)
		return dto.EmptyResponse()
	}
}

func CheckAuthorizationHeader(token string) (*util.Claims, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("authorization header is missing")
	} else if claims, err := util.ParseJwtToken(config.JwtCfg.SigningKey, token); err != nil {
		return nil, fmt.Errorf("token is invalid")
	} else if !infra.RedisClient.Exists(claims.Id) {
		return nil, fmt.Errorf("token is no longer valid")
	} else {
		return claims, nil
	}
}
