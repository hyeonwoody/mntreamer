package model

import (
	"mntreamer/shared/model"
	"time"
)

type Streamer struct {
	PlatformId     uint16 `gorm:"primaryKey;uniqueIndex:idx_streamer"`
	Id             uint32 `gorm:"primaryKey;uniqueIndex:idx_streamer"`
	Nickname       string
	IsStreaming    bool
	Status         int8
	Priority       int8
	Recorded       int16
	LastRecordedAt time.Time
	LastStreamAt   time.Time
}

func (Streamer) TableName() string {
	return "streamer"
}

func NewStreamer(nickname string, platformId uint16) *Streamer {
	return &Streamer{
		PlatformId:     platformId,
		Nickname:       nickname,
		IsStreaming:    false,
		Status:         model.IDLE,
		Priority:       0,
		Recorded:       0,
		LastRecordedAt: time.Now(),
		LastStreamAt:   time.Now(),
	}
}

type ChzzkInformation struct {
	StreamerId  uint16 `gorm:"primaryKey;uniqueIndex:idx_platform"`
	ChannelId   string
	ChannelName string
}
