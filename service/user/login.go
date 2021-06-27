package user

import (
	"fmt"
	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/model"
	"github.com/csumissu/SkyDisk/util/jwt"
	"github.com/gin-gonic/gin"
)

type LoginService struct {
}

func (service *LoginService) Login(c *gin.Context, input dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := model.GetUserByUsername(input.Username)
	if err != nil {
		return nil, fmt.Errorf("username or password is incorrect")
	}
	if authOK, _ := user.CheckPassword(input.Password); !authOK {
		return nil, fmt.Errorf("username or password is incorrect")
	}
	if user.Status == model.Banned {
		return nil, fmt.Errorf("this user was banned")
	}

	claims := jwt.NewClaims()
	claims.Subject = fmt.Sprint(user.ID)
	token, err := jwt.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("internal server error")
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
	return response, nil
}
