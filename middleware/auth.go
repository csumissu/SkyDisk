package middleware

import (
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/service"
	"github.com/csumissu/SkyDisk/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if claims, err := service.CheckAuthorizationHeader(token); err != nil {
			response := dto.Failure(http.StatusUnauthorized, err.Error())
			c.JSON(response.HttpStatus, response)
			c.Abort()
		} else {
			ctx := util.SetCurrentUserID(c.Request.Context(), claims.GetUserID())
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}
	}
}
