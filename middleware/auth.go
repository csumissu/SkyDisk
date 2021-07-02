package middleware

import (
	"context"
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const userIDContextKey = "UserIDContextKey"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if claims, err := user.CheckAuthorizationHeader(token); err != nil {
			response := dto.Failure(http.StatusUnauthorized, err.Error())
			c.JSON(response.HttpStatus, response)
			c.Abort()
		} else {
			subject, _ := strconv.ParseUint(claims.Subject, 10, 32)
			userID := uint(subject)

			ctx := context.WithValue(c.Request.Context(), userIDContextKey, userID)
			c.Request = c.Request.WithContext(ctx)

			c.Set(userIDContextKey, userID)
			c.Next()
		}
	}
}

func GetCurrentUserID(ctx context.Context) uint {
	return ctx.Value(userIDContextKey).(uint)
}
