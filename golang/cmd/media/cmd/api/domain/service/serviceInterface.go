package service

import (
	model "mntreamer/media/cmd/model"
	mntreamerModel "mntreamer/shared/model"
)

type IService interface {
	Download(media *mntreamerModel.Media, channelName string, platformId uint16) error
	Save(platformId uint16, streamId uint32)
	GetFiles(filePath string) ([]model.FileInfo, error)
	GetM3u8(filePath string) ([]model.FileInfo, error)
	GetMediaToRefine() ([]model.MediaRecord, error)
	GetFilePath(mediaRecord *model.MediaRecord, channelName string) string
}
