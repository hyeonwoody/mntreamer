package handler

import (
	api "mntreamer/shared/common/api"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	api.IHandler
	GetFiles(c *gin.Context)
	Stream(c *gin.Context)
	Excise(c *gin.Context)
	GetTargetDuration(c *gin.Context)
}
