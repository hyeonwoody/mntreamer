package repository

import (
	"errors"
	"mntreamer/shared/database"
	mntreamerModel "mntreamer/shared/model"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	mysql database.MysqlWrapper
}

func NewRepository(mysql *database.MysqlWrapper) *Repository {
	return &Repository{mysql: *mysql}
}

func (r *Repository) Create(streamer *mntreamerModel.Streamer) (*mntreamerModel.Streamer, error) {
	result := r.mysql.Driver.Create(streamer)
	if result.Error != nil {
		return nil, result.Error
	}

	return streamer, nil
}

func (r *Repository) Save(target *mntreamerModel.Streamer) (*mntreamerModel.Streamer, error) {
	var savedStreamer *mntreamerModel.Streamer
	err := r.mysql.Driver.Transaction(func(tx *gorm.DB) error {
		result := tx.Save(target)
		if result.Error != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
				existingStreamer := &mntreamerModel.Streamer{}
				err := tx.Where("platform_id = ? AND nickname = ?", target.PlatformId, target.Nickname).First(existingStreamer).Error
				if err != nil {
					return err
				}
				savedStreamer = existingStreamer
				return nil
			}
			return result.Error
		}
		savedStreamer = target
		return nil
	})
	if err != nil {
		return nil, err
	}
	return savedStreamer, nil
}

func (r *Repository) FindByPlatformIdAndStreamerId(platformId uint16, streamerId uint32) (*mntreamerModel.Streamer, error) {
	var streamer mntreamerModel.Streamer
	result := r.mysql.Driver.Where("platform_id = ? AND id = ?", platformId, streamerId).First(&streamer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &streamer, nil
}
