package controllers

import (
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"strconv"
)

func Download(c *gin.Context) {
	objectId, err := strconv.ParseUint(c.Param("objectId"), 10, 32)
	if err == nil {
		response := r.FileService.Download(c, uint(objectId))
		if !cmp.Equal(response, dto.EmptyResponse()) {
			c.JSON(response.HttpStatus, response)
		}
	} else {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
	}
}
