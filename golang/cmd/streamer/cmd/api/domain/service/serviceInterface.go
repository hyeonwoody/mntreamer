package service

import "mntreamer/streamer/cmd/model"

type IService interface {
	GetPlatformIdByName(platform string) (uint16, error)
	SaveStreamer(nickname string, platformId uint16) (*model.Streamer, error)
}
