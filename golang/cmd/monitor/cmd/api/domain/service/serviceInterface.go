package service

import (
	"mntreamer/monitor/cmd/model"

	"gorm.io/gorm"
)

type IService interface {
	AddToMonitoring(platformId uint16, streamerId uint32) (*model.StreamerMonitor, error)
	Checkout() (*model.StreamerMonitor, error)
	UpdateCheckAt(tx *gorm.DB, streamerMonitor *model.StreamerMonitor)
	AddMissCount(monitor *model.StreamerMonitor)
	ResetMissCount(streamerMonitor *model.StreamerMonitor)
}
