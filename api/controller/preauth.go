package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Project-Birdol/birdol-server/controller/jsonmodel"
	"github.com/Project-Birdol/birdol-server/database"
	"github.com/Project-Birdol/birdol-server/database/model"
	"github.com/Project-Birdol/birdol-server/utils/response"
	"github.com/gin-gonic/gin"
)

func ClientVerCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[ClientCheck] ")
		content_type := ctx.GetHeader("Content-Type")
		if content_type != gin.MIMEJSON {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidType)
			ctx.Abort()
			return
		}

		var request_body jsonmodel.ClientVersion
		if err := ctx.ShouldBindJSON(&request_body); err != nil {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrFailParseJSON)
			return
		}

		platform := request_body.Platform
		version := strings.Split(request_body.VersionString, ".")
		build_id := request_body.BuildIdentification

		if len(version) != 3 {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidVersionString)
			return
		}

		var validCli model.ValidClient
		if err := database.Sqldb.Where("platform = ?", platform).First(&validCli).Error; err != nil {
			response.SetErrorResponse(ctx, http.StatusForbidden, response.ErrInvalidPlatform)
			return
		}

		system_ver, _ := strconv.Atoi(version[0])
		major_ver, _ := strconv.Atoi(version[1])
		minor_ver, _ := strconv.Atoi(version[2])
		if (validCli.SystemVersion != uint(system_ver) || validCli.MajorVersion != uint(major_ver) || validCli.MinorVersion != uint(minor_ver)) {
			response.SetErrorResponse(ctx, http.StatusForbidden, response.ErrUpdateRequired)
			return
		}

		if (validCli.Build != build_id) {
			response.SetErrorResponse(ctx, http.StatusForbidden, response.ErrUpdateRequired)
			return
		}

		response.SetNormalResponse(ctx, http.StatusOK, response.ResultOK)
	}
}