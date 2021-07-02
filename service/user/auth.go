package user

import (
	"github.com/csumissu/SkyDisk/config"
	"github.com/csumissu/SkyDisk/infra"
	"github.com/csumissu/SkyDisk/models"
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/util"
	"net/http"
	"strconv"
	"time"
)

type AuthService struct {
}

func (service *AuthService) Login(input dto.LoginRequest) models.ResponseResult {
	user, err := models.GetUserByUsername(input.Username)
	if err != nil {
		return models.Failure(http.StatusBadRequest, "username or password is incorrect")
	}
	if authOK, _ := user.CheckPassword(input.Password); !authOK {
		return models.Failure(http.StatusBadRequest, "username or password is incorrect")
	}
	if user.Status == models.Banned {
		return models.Failure(http.StatusForbidden, "this user is banned")
	}

	expirationDuration := time.Duration(config.JwtCfg.ExpirationHours) * time.Hour
	claims := util.NewClaims(strconv.FormatUint(uint64(user.ID), 10), expirationDuration)
	token, err := util.GenerateJwtToken(config.JwtCfg.SigningKey, claims)
	if err != nil {
		return models.Failure(http.StatusInternalServerError, "token could not be generated")
	}
	infra.Redis().Set(claims.Id, user.ID, expirationDuration)

	response := &dto.LoginResponse{
		UserID:   int(user.ID),
		Nickname: user.Nickname,
		AccessToken: dto.AccessTokenResponse{
			Type:     util.TokenType,
			Token:    token,
			ExpireAt: claims.ExpiresAt * 1000,
		},
	}
	return models.Success(response)
}

func (service *AuthService) Logout(token string) models.ResponseResult {
	if claims, err := util.ParseJwtToken(config.JwtCfg.SigningKey, token); err != nil {
		return models.FailureWithError(http.StatusBadRequest, "token could not be parsed", err)
	} else {
		infra.Redis().Del(claims.Id)
		return models.Success(nil)
	}
}
