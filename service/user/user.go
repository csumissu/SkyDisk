package user

import (
	"github.com/csumissu/SkyDisk/models"
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type UserService struct {
}

func (service *UserService) SearchUserProfile(currentUserId uint) (*dto.UserProfileResponse, error) {
	if user, err := models.GetActiveUserByID(currentUserId); err != nil {
		return nil, gqlerror.Errorf("user could not be found")
	} else {
		return &dto.UserProfileResponse{
			ID:       int(user.ID),
			Username: user.Username,
			Nickname: user.Nickname,
		}, nil
	}
}
