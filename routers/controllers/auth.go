package controllers

import (
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var authService = new(service.AuthService)

func Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err == nil {
		response := authService.Login(request)
		c.JSON(response.HttpStatus, response)
	} else {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
	}
}

func Logout(c *gin.Context) {
	var token = c.GetHeader("Authorization")
	response := authService.Logout(token)
	c.JSON(response.HttpStatus, response)
}
