package user

import (
	"fmt"
	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/model"
	"github.com/csumissu/SkyDisk/util/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthService struct {
}

func (service *AuthService) Login(c *gin.Context, input dto.LoginRequest) model.ResponseResult {
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

	claims := jwt.NewClaims()
	claims.Subject = fmt.Sprint(user.ID)
	token, err := jwt.GenerateToken(claims)
	if err != nil {
		return model.Failure(http.StatusInternalServerError, "token could not be generated")
	}

	response := &dto.LoginResponse{
		UserID:   int(user.ID),
		Nickname: user.Nickname,
		AccessToken: dto.AccessTokenResponse{
			Type:     "Bearer",
			Token:    token,
			ExpireAt: claims.ExpiresAt * 1000,
		},
	}
	return model.Success(response)
}

func (service *AuthService) Logout(c *gin.Context) model.ResponseResult {
	return model.Success(nil)
}
