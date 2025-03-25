package repository

import (
	"mntreamer/media/cmd/model"
	"mntreamer/shared/database"
)

type Repository struct {
	mysql *database.MysqlWrapper
}

func NewRepository(mysql *database.MysqlWrapper) *Repository {
	return &Repository{mysql: mysql}
}

func (r *Repository) Save(mediaRecord *model.MediaRecord) (*model.MediaRecord, error) {
	result := r.mysql.Driver.Model(mediaRecord).Save(mediaRecord)
	if result.Error != nil {
		return nil, result.Error
	}

	return mediaRecord, nil
}
