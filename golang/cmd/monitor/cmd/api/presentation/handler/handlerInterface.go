package handler

import (
	api "mntreamer/shared/common/api"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	api.IHandler
	Add(c *gin.Context)
}
