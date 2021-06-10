package auth

import (
	"context"
	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type LoginService struct {
}

func (service *LoginService) Login(ctx context.Context, input dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := model.GetUserByUsername(input.Username)
	if err != nil {
		return nil, gqlerror.Errorf("username or password is incorrect")
	}
	if authOK, _ := user.CheckPassword(input.Password); !authOK {
		return nil, gqlerror.Errorf("username or password is incorrect")
	}
	if user.Status == model.Banned {
		return nil, gqlerror.Errorf("this user was banned")
	}

	response := &dto.LoginResponse{
		UserID:   int(user.ID),
		Nickname: user.Nickname,
	}
	return response, nil
}
