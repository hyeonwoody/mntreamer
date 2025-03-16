package service

import (
	api "mntreamer/shared/common/api"
	mntreamerModel "mntreamer/shared/model"
)

type IService interface {
	api.IService
	GetPlatformIdByName(name string) (uint16, error)
	BuildStreamer(platformName, nickname string) (*mntreamerModel.Streamer, error)
	GetLiveDetail(streamer *mntreamerModel.Streamer) (*mntreamerModel.Media, error)
}
