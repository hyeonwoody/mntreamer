package handler

import (
	"fmt"
	"mntreamer/monitor/cmd/api/presentation/controller"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctrl controller.IController
}

func NewHandler(ctrl controller.IController) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) Add(c *gin.Context) {
	platform := c.Query("platform")
	nickname := c.Query("nickname")

	if err := h.ctrl.Add(platform, nickname); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("failed to add :%s", nickname))
	}
	c.JSON(http.StatusOK, fmt.Sprintf("successfully added %s", nickname))
}
