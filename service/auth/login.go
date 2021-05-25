package auth

import (
	"context"
	"math/rand"

	"fmt"

	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type LoginService struct {
}

func (service *LoginService) Login(ctx context.Context, input dto.LoginRequest) (*dto.LoginResponse, error) {
	account, err := model.GetAccountByUsername(input.Username)
	if err != nil {
		return nil, gqlerror.Errorf("Username or password is incorrect.")
	}
	if authOK, _ := account.CheckPassword(input.Password); !authOK {
		return nil, gqlerror.Errorf("Username or password is incorrect.")
	}

	response := &dto.LoginResponse{
		UserID:   fmt.Sprintf("T%d", rand.Int()),
		Nickname: "zhangsan",
	}
	return response, nil
}
