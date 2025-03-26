package controller

import (
	model "mntreamer/media/cmd/model"
	api "mntreamer/shared/common/api"
)

type IController interface {
	api.IController
	GetFiles(filePath string) ([]model.FileInfo, error)
	GetFilesToRefine() ([]model.FileInfo, error)
}