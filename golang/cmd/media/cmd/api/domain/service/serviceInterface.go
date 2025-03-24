package service

import (
	mntreamerModel "mntreamer/shared/model"
)

type IService interface {
	Download(media *mntreamerModel.Media, channelName string, platformId uint16) error
	Save(platformId uint16, streamId uint32)
}