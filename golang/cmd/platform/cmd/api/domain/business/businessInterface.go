package business

import (
	mntreamerModel "mntreamer/shared/model"
)

type IBusiness interface {
	GetChannelName(nickname string) (string, error)
	GetChannelId(nickname string) (string, error)
	GetMediaDetail(streamer *mntreamerModel.Streamer) (*mntreamerModel.Media, error)
}

type BusinessStrategy struct {
	clientMap map[uint16]IBusiness
}

func NewBusinessStrategy(clientMap map[uint16]IBusiness) *BusinessStrategy {
	return &BusinessStrategy{clientMap: clientMap}
}

func (strat *BusinessStrategy) GetBusiness(platformId uint16) IBusiness {
	return strat.clientMap[platformId]
}
