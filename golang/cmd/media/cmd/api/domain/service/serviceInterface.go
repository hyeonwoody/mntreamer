package service

import (
	mntreamerModel "mntreamer/shared/model"
)

type IService interface {
	Download(media *mntreamerModel.Media, streamer *mntreamerModel.Streamer) error
}
