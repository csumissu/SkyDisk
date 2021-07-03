package service

import (
	"context"
	"github.com/csumissu/SkyDisk/models"
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/util"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type UserService struct {
}

func (service *UserService) SearchUserProfile(ctx context.Context) (*dto.UserProfileResponse, error) {
	userID := util.GetCurrentUserID(ctx)
	if user, err := models.GetActiveUserByID(userID); err != nil {
		return nil, gqlerror.Errorf("user could not be found")
	} else {
		return &dto.UserProfileResponse{
			ID:       int(user.ID),
			Username: user.Username,
			Nickname: user.Nickname,
		}, nil
	}
}
