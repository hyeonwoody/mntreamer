package service

import (
	"mntreamer/platform/cmd/api/domain/business"
	"mntreamer/platform/cmd/api/infrastructure/repository"
	mntreamerModel "mntreamer/shared/model"
)

type Service struct {
	bizStrat *business.BusinessStrategy
	repo     repository.IRepository
}

func NewService(bizStrat *business.BusinessStrategy, repo repository.IRepository) *Service {
	return &Service{bizStrat: bizStrat, repo: repo}
}

func (s *Service) GetPlatformIdByName(name string) (uint16, error) {
	platform, err := s.repo.FindByName(name)
	if err != nil {
		return 0, err
	}
	return platform.Id, nil
}

func (s *Service) BuildStreamer(platformName, nickname string) (*mntreamerModel.Streamer, error) {
	platformId, err := s.GetPlatformIdByName(platformName)
	if err != nil {
		return nil, err
	}
	clnt := s.bizStrat.GetBusiness(platformId)

	channelName, err := clnt.GetChannelName(nickname)
	if err != nil {
		return nil, err
	}
	channelId, err := clnt.GetChannelId(nickname)
	if err != nil {
		return nil, err
	}
	return mntreamerModel.NewStreamer(nickname, channelName, channelId, platformId), nil
}

func (s *Service) GetLiveDetail(streamer *mntreamerModel.Streamer) (*mntreamerModel.Media, error) {
	clnt := s.bizStrat.GetBusiness(streamer.PlatformId)
	return clnt.GetMediaDetail(streamer)
}
