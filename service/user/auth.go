package user

import (
	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/models"
	"github.com/csumissu/SkyDisk/util/jwt"
	"github.com/csumissu/SkyDisk/util/redis"
	"net/http"
	"strconv"
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

	payload := make(map[string]interface{})
	payload["sub"] = strconv.FormatUint(uint64(user.ID), 10)
	claims := jwt.NewClaims(payload)
	token, err := jwt.GenerateToken(claims)
	if err != nil {
		return models.Failure(http.StatusInternalServerError, "token could not be generated")
	}
	redis.Set(claims["jti"].(string), user.ID, jwt.DefaultExpirationDuration())

	response := &dto.LoginResponse{
		UserID:   int(user.ID),
		Nickname: user.Nickname,
		AccessToken: dto.AccessTokenResponse{
			Type:     jwt.TokenType,
			Token:    token,
			ExpireAt: claims["exp"].(int64) * 1000,
		},
	}
	return models.Success(response)
}

func (service *AuthService) Logout(token string) models.ResponseResult {
	if claims, err := jwt.ParseToken(token); err != nil {
		return models.FailureWithError(http.StatusBadRequest, "token could not be parsed", err)
	} else {
		redis.Del(claims["jti"].(string))
		return models.Success(nil)
	}
}
