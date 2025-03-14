package repository

import (
	"mntreamer/monitor/cmd/model"
	"mntreamer/shared/database"
)

type Repository struct {
	mysql database.MysqlWrapper
}

func NewRepository(mysql *database.MysqlWrapper) *Repository {
	return &Repository{mysql: *mysql}
}

func (r *Repository) Save(streamerMonitor *model.StreamerMonitor) (*model.StreamerMonitor, error) {

	result := r.mysql.Driver.Create(streamerMonitor)
	if result.Error != nil {
		return nil, result.Error
	}

	return streamerMonitor, nil
}
