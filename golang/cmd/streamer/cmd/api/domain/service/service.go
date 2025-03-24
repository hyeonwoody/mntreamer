package service

import (
	mntreamerModel "mntreamer/shared/model"
	"mntreamer/streamer/cmd/api/infrastructure/repository"
	"time"
)

type Service struct {
	streamerRepo repository.IRepository
}

func NewService(repo repository.IRepository) *Service {
	return &Service{streamerRepo: repo}
}

func (s *Service) GetPlatformIdByName(platform string) (uint16, error) {
	return 0, nil
}

func (s *Service) Create(streamer *mntreamerModel.Streamer) (*mntreamerModel.Streamer, error) {
	return s.streamerRepo.Create(streamer)
}

func (s *Service) Save(streamer *mntreamerModel.Streamer) (*mntreamerModel.Streamer, error) {
	return s.streamerRepo.Save(streamer)
}

func (s *Service) FindByPlatformIdAndStreamerId(platformId uint16, streamerId uint32) (*mntreamerModel.Streamer, error) {
	return s.streamerRepo.FindByPlatformIdAndStreamerId(platformId, streamerId)
}

func (s *Service) CheckMonitoringEligibility(streamer *mntreamerModel.Streamer) bool {
	if streamer.Status == mntreamerModel.IDLE {
		return true
	}
	return false
}

func (s *Service) UpdateStatus(streamer *mntreamerModel.Streamer, status int8) {
	streamer.Status = status
	s.streamerRepo.Save(streamer)
}

func (s *Service) UpdateStatusWithId(platformId uint16, streamerId uint32, status int8) {
	streamer, err := s.streamerRepo.FindByPlatformIdAndStreamerId(platformId, streamerId)
	if err != nil {
		return
	}
	streamer.Status = status
	s.streamerRepo.Save(streamer)
}

func (s *Service) UpdateLastRecordedAt(streamer *mntreamerModel.Streamer) {
	streamer.LastRecordedAt = time.Now()
	s.streamerRepo.Save(streamer)
}

func (s *Service) UpdateLastStreamAt(streamer *mntreamerModel.Streamer) {
	streamer.LastStreamAt = time.Now()
	s.streamerRepo.Save(streamer)
}
