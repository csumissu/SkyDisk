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
			response := dto.Failure(http.StatusUnauthorized, err.Error())
			c.JSON(response.HttpStatus, response)
			c.Abort()
		} else {
			userID, _ := strconv.ParseUint(claims.Subject, 10, 32)
			c.Set(userIDContextKey, userID)
			c.Next()
		}
	}
}
