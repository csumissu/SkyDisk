package middleware

import (
	"github.com/csumissu/SkyDisk/config"
	"github.com/csumissu/SkyDisk/models"
	"github.com/csumissu/SkyDisk/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const userIDContextKey = "UserIDContextKey"

func SignRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := c.GetHeader("Authorization"); len(token) == 0 {
			response := models.Failure(http.StatusUnauthorized, "authorization header is missing")
			c.JSON(response.HttpStatus, response)
			c.Abort()
		} else if claims, err := util.ParseJwtToken(config.JwtCfg.SigningKey, token); err != nil {
			response := models.FailureWithError(http.StatusUnauthorized, "token is invalid", err)
			c.JSON(response.HttpStatus, response)
			c.Abort()
		} else if userID, err := strconv.ParseUint(claims.Subject, 10, 32); err != nil {
			response := models.FailureWithError(http.StatusUnauthorized, "token is invalid", err)
			c.JSON(response.HttpStatus, response)
			c.Abort()
		} else {
			c.Set(userIDContextKey, userID)
			c.Next()
		}
	}
}
