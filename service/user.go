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
	if user, err := GetCurrentUser(ctx); err != nil {
		return nil, err
	} else {
		return &dto.UserProfileResponse{
			ID:       int(user.ID),
			Username: user.Username,
			Nickname: user.Nickname,
		}, nil
	}
}

func GetCurrentUser(ctx context.Context) (*models.User, error) {
	if userID, ok := util.GetCurrentUserID(ctx); ok {
		if user, err := models.GetActiveUserByID(userID); err == nil {
			return user, nil
		}
	}
	return nil, gqlerror.Errorf("user could not be found")
}
