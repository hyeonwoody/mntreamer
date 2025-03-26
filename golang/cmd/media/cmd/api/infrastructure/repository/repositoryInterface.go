package repository

import (
	model "mntreamer/media/cmd/model"
)

type IRepository interface {
	Save(mediaRecord *model.MediaRecord) (*model.MediaRecord, error)
	FindByStatus(status int) ([]model.MediaRecord, error)
}