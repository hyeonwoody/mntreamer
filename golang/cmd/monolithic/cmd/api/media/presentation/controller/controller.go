package controller

import (
	media "mntreamer/media/cmd/api/domain/service"
	"mntreamer/media/cmd/model"
)

type ControllerMono struct {
	svc media.IService
}

func NewControllerMono(svc media.IService) *ControllerMono {
	return &ControllerMono{svc: svc}
}

func (c *ControllerMono) GetFiles(filePath string) ([]model.FileInfo, error) {
	return c.svc.GetFiles(filePath)
}
