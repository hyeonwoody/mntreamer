package repository

import (
	model "mntreamer/media/cmd/model"
	"time"
)

type IRepository interface {
	Terminate(platformId uint16, streamerId uint32, date time.Time, sequence uint16) (*model.MediaRecord, error)
	Save(mediaRecord *model.MediaRecord) (*model.MediaRecord, error)
	FindByStatus(status int) ([]model.MediaRecord, error)
}
