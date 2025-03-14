package service

import "mntreamer/monitor/cmd/model"

type IService interface {
	AddToMonitoring(platformId uint16, streamerId uint32) (*model.StreamerMonitor, error)
}
