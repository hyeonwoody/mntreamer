package service

import (
	mntreamerModel "mntreamer/shared/model"
)

type IService interface {
	GetPlatformIdByName(platform string) (uint16, error)

	Create(streamer *mntreamerModel.Streamer) (*mntreamerModel.Streamer, error)
	Save(streamer *mntreamerModel.Streamer) (*mntreamerModel.Streamer, error)
	FindByPlatformIdAndStreamerId(platformId uint16, streamerId uint32) (*mntreamerModel.Streamer, error)
	CheckMonitoringEligibility(streamer *mntreamerModel.Streamer) bool
	UpdateStatus(streamer *mntreamerModel.Streamer, status int8)
	UpdateStatusWithId(platforId uint16, streamerId uint32, status int8)
	UpdateLastRecordedAt(streamer *mntreamerModel.Streamer)
	UpdateLastStreamAt(streamer *mntreamerModel.Streamer)
}
