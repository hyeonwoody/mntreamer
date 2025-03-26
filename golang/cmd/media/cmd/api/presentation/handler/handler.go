package handler

import (
	"fmt"
	"mntreamer/media/cmd/api/presentation/controller"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctrl     controller.IController
	basePath string
}

func NewHandler(basePath string, ctrl controller.IController) *Handler {
	return &Handler{basePath: basePath, ctrl: ctrl}
}

func (h *Handler) GetFiles(c *gin.Context) {
	var req struct {
		Path string `json:"path"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fileInfos interface{}
	var err error

	switch req.Path {
	case "MEDIARECORD":
		fileInfos, err = h.ctrl.GetFilesToRefine()
	default:
		if !strings.HasPrefix(req.Path, h.basePath) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Cannot access path outside base directory: %s", req.Path)})
			return
		}
		fileInfos, err = h.ctrl.GetFiles(req.Path)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve files: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": fileInfos})
}
