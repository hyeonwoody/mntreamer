package service

import (
	model "mntreamer/media/cmd/model"
	mntreamerModel "mntreamer/shared/model"
	"os"
	"time"
)

type IService interface {
	GetRootPath() string
	Download(media *mntreamerModel.Media, channelName string, platformId uint16) error
	Save(platformId uint16, streamerId uint32, status int8)
	GetFilePath(record *model.MediaRecord, platformName string, name string) string
	GetFiles(filePath string) ([]model.FileInfo, error)
	GetM3u8(filePath string, sequence uint16) ([]model.FileInfo, error)
	GetMediaToRefine() ([]model.MediaRecord, error)
	GetPlatformNameByFilePath(filePath string) (string, error)
	GetChannelNameByFilePath(filePath string) (string, error)
	GetDateByFilePath(fullPath string) (time.Time, error)
	GetSequenceByFilePath(fullPath string) (uint16, error)
	Stream(filePath string) (string, error)
	StreamMediaPlaylist(filePath string) (string, error)
	StreamSegment(filePath string) (string, error)
	StreamMp4(filePath string) (*os.File, error)
	Decode(filePath string) (interface{}, error)
	Excise(path string, begin float64, end float64) error
	UpdateStatus(mediaRecord *model.MediaRecord, status int8) (*model.MediaRecord, error)
	Confirm(platformId uint16, streamerId uint32, fullPath string) (*model.MediaRecord, error)
	Delete(platformId uint16, streamerId uint32, fullPath string) (*model.MediaRecord, error)
}
