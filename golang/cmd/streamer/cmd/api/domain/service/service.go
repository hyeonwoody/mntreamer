package service

import (
	"mntreamer/streamer/cmd/api/infrastructure/repository"
	"mntreamer/streamer/cmd/model"
)

type Service struct {
	repo repository.IRepository
}

func NewService(repo repository.IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetPlatformIdByName(platform string) (uint16, error) {
	return 0, nil
}

func (s *Service) SaveStreamer(nickname string, platformId uint16) (*model.Streamer, error) {
	return s.repo.Save(model.NewStreamer(nickname, platformId))
}
