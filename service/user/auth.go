package user

import (
	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/model"
	"github.com/csumissu/SkyDisk/util/jwt"
	"github.com/csumissu/SkyDisk/util/redis"
	"net/http"
)

type AuthService struct {
}

func (service *AuthService) Login(input dto.LoginRequest) model.ResponseResult {
	user, err := model.GetUserByUsername(input.Username)
	if err != nil {
		return model.Failure(http.StatusBadRequest, "username or password is incorrect")
	}
	if authOK, _ := user.CheckPassword(input.Password); !authOK {
		return model.Failure(http.StatusBadRequest, "username or password is incorrect")
	}
	if user.Status == model.Banned {
		return model.Failure(http.StatusForbidden, "this user is banned")
	}

	payload := make(map[string]interface{})
	payload["sub"] = user.ID
	claims := jwt.NewClaims(payload)
	token, err := jwt.GenerateToken(claims)
	if err != nil {
		return model.Failure(http.StatusInternalServerError, "token could not be generated")
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
	return model.Success(response)
}

func (service *AuthService) Logout(token string) model.ResponseResult {
	if claims, err := jwt.ParseToken(token); err != nil {
		return model.FailureWithError(http.StatusBadRequest, "token could not be parsed", err)
	} else {
		redis.Del(claims["jti"].(string))
		return model.Success(nil)
	}
}
