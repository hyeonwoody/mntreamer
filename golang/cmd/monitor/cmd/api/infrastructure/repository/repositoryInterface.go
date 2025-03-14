package repository

import (
	model "mntreamer/monitor/cmd/model"
)

type IRepository interface {
	Save(streamerMonitor *model.StreamerMonitor) (*model.StreamerMonitor, error)
}
