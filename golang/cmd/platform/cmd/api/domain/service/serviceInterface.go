package service

import (
	api "mntreamer/shared/common/api"
)

type IService interface {
	api.IService
	GetPlatformIdByName(name string) (uint16, error)
}
