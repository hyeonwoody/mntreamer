package repository

import (
	"mntreamer/monitor/cmd/model"
	"mntreamer/shared/database"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	mysql database.MysqlWrapper
}

func NewRepository(mysql *database.MysqlWrapper) *Repository {
	return &Repository{mysql: *mysql}
}

func (r *Repository) Create(streamerMonitor *model.StreamerMonitor) (*model.StreamerMonitor, error) {

	result := r.mysql.Driver.Create(streamerMonitor)
	if result.Error != nil {
		return nil, result.Error
	}

	return streamerMonitor, nil
}

func (r *Repository) Save(streamerMonitor *model.StreamerMonitor) (*model.StreamerMonitor, error) {
	result := r.mysql.Driver.Model(streamerMonitor).Save(streamerMonitor)
	if result.Error != nil {
		return nil, result.Error
	}

	return streamerMonitor, nil
}

func (r *Repository) FindByCheckAtLock(currentTime time.Time) (*model.StreamerMonitor, *gorm.DB, error) {
	var streamerMonitor model.StreamerMonitor
	tx := r.mysql.Driver.Begin()
	result := tx.Clauses(clause.Locking{
		Strength: "UPDATE",
		Table:    clause.Table{Name: "streamer_monitor"},
	}).
		Where("check_at < ?", currentTime).
		First(&streamerMonitor)
	if result.Error != nil {
		return nil, tx, result.Error
	}
	return &streamerMonitor, tx, nil
}

func (r *Repository) UpdateTx(tx *gorm.DB, streamerMonitor *model.StreamerMonitor) {
	result := tx.Model(streamerMonitor).Save(streamerMonitor)
	if result.Error != nil {
		tx.Rollback()
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return
	}
}
