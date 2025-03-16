package repository

import (
	model "mntreamer/monitor/cmd/model"
	"time"

	"gorm.io/gorm"
)

type IRepository interface {
	Save(streamerMonitor *model.StreamerMonitor) (*model.StreamerMonitor, error)
	Create(streamerMonitor *model.StreamerMonitor) (*model.StreamerMonitor, error)
	FindByCheckAtLock(time time.Time) (*model.StreamerMonitor, *gorm.DB, error)
	UpdateTx(tx *gorm.DB, streamerMonitor *model.StreamerMonitor)
}
