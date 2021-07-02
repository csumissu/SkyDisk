package middleware

import (
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const userIDContextKey = "UserIDContextKey"

func SignRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if claims, err := user.CheckAuthorizationHeader(token); err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse(err))
			c.Abort()
		} else {
			userID, _ := strconv.ParseUint(claims.Subject, 10, 32)
			c.Set(userIDContextKey, userID)
			c.Next()
		}
	}
}
