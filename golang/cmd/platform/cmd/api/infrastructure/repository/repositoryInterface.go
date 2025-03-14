package repository

import (
	"mntreamer/platform/cmd/model"
	api "mntreamer/shared/common/api"
)

type IRepository interface {
	api.IRepository
	FindByName(name string) (*model.Platform, error)
}
