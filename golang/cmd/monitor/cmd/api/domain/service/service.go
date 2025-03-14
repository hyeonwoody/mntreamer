package service

import (
	"mntreamer/monitor/cmd/api/infrastructure/repository"
	"mntreamer/monitor/cmd/model"
)

type Service struct {
	repo repository.IRepository
}

func NewService(repo repository.IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddToMonitoring(platformId uint16, streamerId uint32) (*model.StreamerMonitor, error) {
	return s.repo.Save(model.NewStreamerMonitor(platformId, streamerId))
}
