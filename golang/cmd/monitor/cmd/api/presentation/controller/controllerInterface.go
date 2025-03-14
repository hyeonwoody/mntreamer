package controller

import (
	api "mntreamer/shared/common/api"
)

type IController interface {
	api.IController
	Add(platform, nickname string) error
}
