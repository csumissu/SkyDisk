package auth

import (
	"context"
	"math/rand"

	"fmt"

	"github.com/csumissu/SkyDisk/graph/model"
)

type LoginService struct {
}

func (service *LoginService) Login(ctx context.Context, input model.LoginRequest) (*model.LoginResponse, error) {
	response := &model.LoginResponse{
		ID:       fmt.Sprintf("T%d", rand.Int()),
		Nickname: "zhangsan",
	}
	return response, nil
}
