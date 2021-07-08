package controllers

import (
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"strconv"
)

func DownloadObject(c *gin.Context) {
	objectID, err := strconv.ParseUint(c.Param("objectID"), 10, 32)
	if err == nil {
		response := r.FileService.DownloadObject(c, uint(objectID))
		if !cmp.Equal(response, dto.EmptyResponse()) {
			c.JSON(response.HttpStatus, response)
		}
	} else {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
	}
}
