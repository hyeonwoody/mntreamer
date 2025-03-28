package controller

import (
	model "mntreamer/media/cmd/model"
	api "mntreamer/shared/common/api"
)

type IController interface {
	api.IController
	GetFiles(filePath string) ([]model.FileInfo, error)
	GetTargetDuration(path string) (float64, error)
	GetFilesToRefine() ([]model.FileInfo, error)
	Stream(filePath string) (string, error)
	Excise(path string, begin float64, end float64) error
	Confirm(filePath string) error
	Delete(filePath string) error
}
