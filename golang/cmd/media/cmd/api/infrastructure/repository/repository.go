package repository

import (
	"mntreamer/media/cmd/model"
	"mntreamer/shared/database"
	"time"
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

func (r *Repository) Terminate(platformId uint16, streamerId uint32, date time.Time, sequence uint16) (*model.MediaRecord, error) {
	var mediaRecord model.MediaRecord

	result := r.mysql.Driver.Where("platform_id = ? AND streamer_id = ? AND date = ? AND sequence = ?", platformId, streamerId, date, sequence).
		First(&mediaRecord)

	if result.Error != nil {
		return nil, result.Error
	}

	result = r.mysql.Driver.Model(&mediaRecord).Update("status", 8) // Assuming 2 is the status for terminated
	if result.Error != nil {
		return nil, result.Error
	}
	mediaRecord.Status = 8
	return &mediaRecord, nil
}

func (r *Repository) FindByStatus(status int8) ([]model.MediaRecord, error) {
	var mediaRecords []model.MediaRecord
	result := r.mysql.Driver.Where("status = ?", status).
		Find(&mediaRecords)
	if result.Error != nil {
		return nil, result.Error
	}
	return mediaRecords, nil
}
