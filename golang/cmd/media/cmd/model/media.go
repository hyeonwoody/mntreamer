package model

import (
	"time"
)

type MediaRecord struct {
	PlatformId uint16 `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_platform_streamer_date_sequence"`
	StreamerId uint32 `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_platform_streamer_date_sequence"`
	Status     int8
	Date       time.Time `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_platform_streamer_date_sequence"`
	Sequence   uint16    `gorm:"primaryKey;autoIncrement:true;uniqueIndex:idx_platform_streamer_date_sequence"`
}

func (MediaRecord) TableName() string {
	return "media_record"
}

func NewMediaRecord(platformId uint16, streamerId uint32, status int8) *MediaRecord {
	now := time.Now()
	year, month, day := now.Date()
	return &MediaRecord{
		PlatformId: platformId,
		StreamerId: streamerId,
		Status:     status,
		Date:       time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
	}
}

func NewInstance(platformId uint16, streamerId uint32, date time.Time, sequence uint16, status int8) *MediaRecord {
	return &MediaRecord{
		PlatformId: platformId,
		StreamerId: streamerId,
		Date:       date,
		Sequence:   sequence,
		Status:     status,
	}
}
