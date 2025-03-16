package repository

import (
	"mntreamer/platform/cmd/model"
	database "mntreamer/shared/database"
)

type Repository struct {
	mysql *database.MysqlWrapper
}

func NewRepository(mysql *database.MysqlWrapper) *Repository {
	return &Repository{mysql: mysql}
}

func (r *Repository) FindByName(name string) (*model.Platform, error) {
	var platform model.Platform
	result := r.mysql.Driver.Table("platform").
		Where("name = ?", name).Find(&platform)
	if result.Error != nil {
		return nil, result.Error
	}
	return &platform, nil
}
