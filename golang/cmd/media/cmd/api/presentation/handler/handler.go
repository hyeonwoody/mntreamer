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
	rootPath string
}

func NewHandler(rootPath string, ctrl controller.IController) *Handler {
	return &Handler{rootPath: rootPath, ctrl: ctrl}
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
		if !strings.HasPrefix(req.Path, h.rootPath) {
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

func (h *Handler) GetTargetDuration(c *gin.Context) {
	filePath := c.Param("filePath")
	targetDuration, err := h.ctrl.GetTargetDuration(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"targetDuration": targetDuration})
}

func (h *Handler) Stream(c *gin.Context) {
	filePath := c.Param("filePath")

	if fullPath, err := h.ctrl.Stream(filePath); err == nil {
		c.Header("Content-Type", "application/vnd.apple.mpegurl")
		c.File(fullPath)
		return
	}
	c.String(http.StatusNotFound, "File not found")
}

func (h *Handler) Excise(c *gin.Context) {
	var req struct {
		FullPath string  `json:"fullPath"`
		Begin    float64 `json:"begin"`
		End      float64 `json:"end"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ctrl.Excise(req.FullPath, req.Begin, req.End); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) Confirm(c *gin.Context) {
	filePath := c.Param("filePath")
	if err := h.ctrl.Confirm(filePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	filePath := c.Param("filePath")
	if err := h.ctrl.Delete(filePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusOK)
}
