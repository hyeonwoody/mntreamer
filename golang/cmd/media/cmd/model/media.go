package model

import "time"

type MediaRecord struct {
	PlatformId uint16    `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_mediaRecord_platforId"`
	StreamerId uint32    `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_mediaRecord_streamerId"`
	Date       time.Time `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_mediaRecord_date"`
}

func (MediaRecord) TableName() string {
	return "media_record"
}

func NewMediaRecord(platformId uint16, streamerId uint32) *MediaRecord {
	now := time.Now()
	year, month, day := now.Date()
	return &MediaRecord{
		PlatformId: platformId,
		StreamerId: streamerId,
		Date:       time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
	}
}
