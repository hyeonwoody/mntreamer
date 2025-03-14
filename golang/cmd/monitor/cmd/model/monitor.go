package model

import "time"

type StreamerMonitor struct {
	Platform_id uint16
	Streamer_id uint32
	CheckAt     time.Time
}

func (StreamerMonitor) TableName() string {
	return "streamer_monitor"
}

func NewStreamerMonitor(platforId uint16, streamerId uint32) *StreamerMonitor {
	return &StreamerMonitor{
		Platform_id: platforId,
		Streamer_id: streamerId,
		CheckAt:     time.Now(),
	}
}
