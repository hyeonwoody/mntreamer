package handler

import (
	"fmt"
	"mntreamer/media/cmd/api/presentation/controller"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctrl controller.IController
}

func NewHandler(ctrl controller.IController) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetFiles(c *gin.Context) {
	var req struct {
		Path string `json:"path"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetPath := req.Path
	if targetPath == "" {
		targetPath = "/"
	}
	fileInfos, err := h.ctrl.GetFiles(targetPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("failed to find path :%s", targetPath))
	}
	c.JSON(http.StatusOK, gin.H{"files": fileInfos})
}
