package middleware

import (
	"github.com/csumissu/SkyDisk/model"
	"github.com/csumissu/SkyDisk/util/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const userIDContextKey = "UserIDContextKey"

func SignRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := c.GetHeader("Authorization"); len(token) == 0 {
			response := model.Failure(http.StatusUnauthorized, "authorization header is missing")
			c.JSON(response.HttpStatus, response)
			c.Abort()
		} else if claims, err := jwt.ParseToken(token); err != nil {
			response := model.FailureWithError(http.StatusUnauthorized, "token is invalid", err)
			c.JSON(response.HttpStatus, response)
			c.Abort()
		} else if userID, err := strconv.ParseUint(claims["sub"].(string), 10, 32); err != nil {
			response := model.FailureWithError(http.StatusUnauthorized, "token is invalid", err)
			c.JSON(response.HttpStatus, response)
			c.Abort()
		} else {
			c.Set(userIDContextKey, userID)
			c.Next()
		}
	}
}
