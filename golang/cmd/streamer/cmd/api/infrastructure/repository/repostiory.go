package repository

import (
	"mntreamer/shared/database"
	"mntreamer/streamer/cmd/model"
)

type Repository struct {
	mysql database.MysqlWrapper
}

func NewRepository(mysql *database.MysqlWrapper) *Repository {
	return &Repository{mysql: *mysql}
}

func (r *Repository) Save(target *model.Streamer) (*model.Streamer, error) {

	result := r.mysql.Driver.Create(target)
	if result.Error != nil {
		return nil, result.Error
	}

	return target, nil
}
