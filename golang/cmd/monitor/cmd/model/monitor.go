package model

import "time"

type StreamerMonitor struct {
	PlatformId uint16 `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_platform_streamerId"`
	StreamerId uint32 `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_platform_streamerId"`
	CheckAt    time.Time
	MissCount  uint8
}

func (StreamerMonitor) TableName() string {
	return "streamer_monitor"
}

func NewStreamerMonitor(platforId uint16, streamerId uint32) *StreamerMonitor {
	return &StreamerMonitor{
		PlatformId: platforId,
		StreamerId: streamerId,
		CheckAt:    time.Now(),
		MissCount:  0,
	}
}
